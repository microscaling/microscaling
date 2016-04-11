// Schedule using docker remote API
package docker

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fsouza/go-dockerclient"
	"github.com/op/go-logging"

	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/scheduler"
)

const labelMap string = "com.microscaling.microscaling-in-a-box"

var log = logging.MustGetLogger("mssscheduler")

type dockerContainer struct {
	state   string
	updated bool
}

type DockerScheduler struct {
	client         *docker.Client
	pullImages     bool
	taskContainers map[string]map[string]*dockerContainer // tasks indexed by app name, containers indexed by ID
	sync.Mutex
}

func NewScheduler(pullImages bool, dockerHost string) *DockerScheduler {
	client, err := docker.NewClient(dockerHost)
	if err != nil {
		log.Errorf("Error starting Docker client: %v", err)
		return nil
	}

	return &DockerScheduler{
		client:         client,
		taskContainers: make(map[string]map[string]*dockerContainer),
		pullImages:     pullImages,
	}
}

// compile-time assert that we implement the right interface
var _ scheduler.Scheduler = (*DockerScheduler)(nil)

var scaling sync.WaitGroup

func (c *DockerScheduler) InitScheduler(appId string, task *demand.Task) (err error) {
	log.Infof("Docker initializing task %s", appId)

	c.Lock()
	defer c.Unlock()

	c.taskContainers[appId] = make(map[string]*dockerContainer, 100)

	// We may need to pull the image for this container
	if c.pullImages {
		pullOpts := docker.PullImageOptions{
			Repository: task.Image,
		}

		authOpts := docker.AuthConfiguration{}

		log.Infof("Pulling image: %v", task.Image)
		err = c.client.PullImage(pullOpts, authOpts)
		if err != nil {
			log.Errorf("Failed to pull image %s: %v", task.Image, err)
		}
	}

	return err
}

// startTask creates the container and then starts it
func (c *DockerScheduler) startTask(name string, task *demand.Task) {
	var labels map[string]string = map[string]string{
		labelMap: name,
	}

	var cmds []string = strings.Fields(task.Command)

	createOpts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:        task.Image,
			Cmd:          cmds,
			AttachStdout: true,
			AttachStdin:  true,
			Labels:       labels,
		},
	}

	hostConfig := docker.HostConfig{
		PublishAllPorts: task.PublishAllPorts,
	}

	go func() {
		scaling.Add(1)
		defer scaling.Done()

		log.Debugf("[start] task %s", name)
		container, err := c.client.CreateContainer(createOpts)
		if err != nil {
			log.Errorf("Couldn't create container for task %s: %v", name, err)
			return
		}

		var containerID = container.ID[:12]

		c.Lock()
		c.taskContainers[name][containerID] = &dockerContainer{
			state: "created",
		}
		c.Unlock()
		log.Debugf("[created] task %s ID %s", name, containerID)

		// Start it
		err = c.client.StartContainer(containerID, &hostConfig)
		if err != nil {
			log.Errorf("Couldn't start container ID %s for task %s: %v", containerID, name, err)
			return
		}

		log.Debugf("[starting] task %s ID %s", name, containerID)

		c.Lock()
		c.taskContainers[name][containerID].state = "starting"
		c.Unlock()
	}()
}

// stopTask kills the last container we know about of this type
func (c *DockerScheduler) stopTask(name string, task *demand.Task) error {
	var err error = nil

	// Kill a currently-running container of this type
	c.Lock()
	theseContainers := c.taskContainers[name]
	var containerToKill string
	for id, v := range theseContainers {
		if v.state == "running" {
			containerToKill = id
			v.state = "stopping"
			break
		}
	}
	c.Unlock()

	if containerToKill == "" {
		return fmt.Errorf("[stop] No containers of type %s to kill", name)
	}

	removeOpts := docker.RemoveContainerOptions{
		ID:            containerToKill,
		RemoveVolumes: true,
	}

	go func() {
		scaling.Add(1)
		defer scaling.Done()

		log.Debugf("[stopping] container for task %s with ID %s", name, containerToKill)
		err = c.client.StopContainer(containerToKill, 1)
		if err != nil {
			log.Errorf("Couldn't stop container %s: %v", containerToKill, err)
			return
		}

		c.Lock()
		c.taskContainers[name][containerToKill].state = "removing"
		c.Unlock()

		log.Debugf("[removing] container for task %s with ID %s", name, containerToKill)
		err = c.client.RemoveContainer(removeOpts)
		if err != nil {
			log.Errorf("Couldn't remove container %s: %v", containerToKill, err)
			return
		}
	}()

	return nil
}

func (c *DockerScheduler) StopStartTasks(tasks map[string]demand.Task) error {
	// Create containers if there aren't enough of them, and stop them if there are too many
	var too_many []string
	var too_few []string
	var diff int
	var err error = nil

	// TODO: Consider checking the number running before we start & stop
	for name, task := range tasks {
		if task.Demand > task.Requested {
			// There aren't enough of these containers yet
			too_few = append(too_few, name)
		}

		if task.Demand < task.Requested {
			// There aren't enough of these containers yet
			too_many = append(too_many, name)
		}
	}

	// Scale down first to free up resources
	for _, name := range too_many {
		task := tasks[name]
		diff = task.Requested - task.Demand
		log.Debugf("Stop %d of task %s", diff, name)
		for i := 0; i < diff; i++ {
			err = c.stopTask(name, &task)
			if err != nil {
				log.Errorf("Couldn't stop %s: %v ", name, err)
			}
			task.Requested -= 1
		}
		tasks[name] = task
	}

	// Now we can scale up
	for _, name := range too_few {
		task := tasks[name]
		diff = task.Demand - task.Requested
		log.Debugf("Start %d of task %s", diff, name)
		for i := 0; i < diff; i++ {
			c.startTask(name, &task)
			task.Requested += 1
		}
		tasks[name] = task
	}

	// Don't return until all the scale tasks are complete
	scaling.Wait()
	return err
}

func statusToState(status string) string {
	if strings.Contains(status, "Up") {
		return "running"
	}
	if strings.Contains(status, "Removal") {
		return "removing"
	}
	if strings.Contains(status, "Exit") {
		return "exited"
	}
	if strings.Contains(status, "Dead") {
		return "dead"
	}
	log.Errorf("Unexpected docker status %s", status)
	return "unknown"
}

func (c *DockerScheduler) CountAllTasks(running *demand.Tasks) error {
	// Docker Remote API https://docs.docker.com/reference/api/docker_remote_api_v1.20/
	// get /containers/json
	var err error
	var containers []docker.APIContainers
	containers, err = c.client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return fmt.Errorf("Failed to list containers: %v", err)
	}

	running.Lock()
	defer running.Unlock()
	c.Lock()
	defer c.Unlock()

	// Reset all the running counts to 0
	tasks := running.Tasks
	for name, t := range tasks {
		t.Running = 0
		tasks[name] = t

		for _, cc := range c.taskContainers[name] {
			cc.updated = false
		}
	}

	var taskName string
	var present bool

	for i := range containers {
		labels := containers[i].Labels
		taskName, present = labels[labelMap]
		if present {
			// Only update tasks that are already in our task map - don't try to manage anything else
			// log.Debugf("Found a container with labels %v", labels)
			t, inOurTasks := tasks[taskName]
			if inOurTasks {
				newState := statusToState(containers[i].Status)
				id := containers[i].ID[:12]
				thisContainer, ok := c.taskContainers[taskName][id]
				if !ok {
					log.Infof("We have no previous record of container %s, state %s", id, newState)
					thisContainer = &dockerContainer{}
					c.taskContainers[taskName][id] = thisContainer
				}

				switch newState {
				case "running":
					t.Running++
					// We could be moving from starting to running, or it could be a container that's totally new to us
					if thisContainer.state == "starting" || thisContainer.state == "" {
						thisContainer.state = newState
					}
				case "removing":
					if thisContainer.state != "removing" {
						log.Errorf("Container %s is being removed, but we didn't terminate it", id)
					}
				case "exited":
					if thisContainer.state != "stopping" && thisContainer.state != "exited" {
						log.Errorf("Container %s is being removed, but we didn't terminate it", id)
					}
				case "dead":
					if thisContainer.state != "dead" {
						log.Errorf("Container %s is dead", id)
					}
					thisContainer.state = newState
				}

				thisContainer.updated = true
				tasks[taskName] = t
			}
		}
	}

	for name, task := range tasks {
		log.Debugf("  %s: internally running %d, requested %d", name, task.Running, task.Requested)
		for id, cc := range c.taskContainers[name] {
			log.Debugf("  %s - %s", id, cc.state)
			if !cc.updated {
				if cc.state == "removing" || cc.state == "exited" {
					log.Debugf("    Deleting %s", id)
					delete(c.taskContainers[name], id)
				} else if cc.state != "created" && cc.state != "starting" && cc.state != "stopping" {
					log.Errorf("Bad state for container %s: %s", id, cc.state)
				}
			}
		}
	}

	return err
}
