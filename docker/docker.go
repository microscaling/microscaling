// Schedule using docker compose
package docker

import (
	"fmt"
	"log"
	"os"

	"github.com/force12io/force12/demand"
	"github.com/force12io/force12/scheduler"

	"github.com/fsouza/go-dockerclient"
)

const f12_map string = "io.force12.microscaling-in-a-box"

type DockerScheduler struct {
	client     *docker.Client
	hostConfig docker.HostConfig
	containers map[string][]string
}

func NewScheduler() *DockerScheduler {
	client, _ := docker.NewClient(os.Getenv("DOCKER_HOST"))

	return &DockerScheduler{
		client:     client,
		containers: make(map[string][]string),
	}
}

// compile-time assert that we implement the right interface
var _ scheduler.Scheduler = (*DockerScheduler)(nil)

func (c *DockerScheduler) InitScheduler(appId string, task *demand.Task) error {
	log.Printf("Docker initializing task %s", appId)

	// TODO Size of this should be max containers
	c.containers[appId] = make([]string, 100)
	return nil
}

// startTask creates the container and then starts it
func (c *DockerScheduler) startTask(name string, task *demand.Task) error {
	var err error = nil
	var labels map[string]string = map[string]string{
		f12_map: name,
	}

	createOpts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:        task.Image,
			AttachStdout: true,
			AttachStdin:  true,
			Labels:       labels,
		},
	}

	container, err := c.client.CreateContainer(createOpts)
	if err != nil {
		return err
	}

	c.containers[name] = append(c.containers[name], container.ID[:12])
	log.Printf("Created task %s with image %s, ID %s", name, task.Image, container.ID)

	// Start it
	err = c.client.StartContainer(container.ID, &c.hostConfig)

	return err
}

// stopTask kills the last container we know about of this type
func (c *DockerScheduler) stopTask(name string, task *demand.Task) error {
	var err error = nil

	// Kill the last container of this type.
	these_containers := c.containers[name]
	container_to_kill := these_containers[len(these_containers)-1]
	c.containers[name] = these_containers[:len(these_containers)-1]
	log.Printf("Removing task %s with ID %s", name, container_to_kill)

	err = c.client.StopContainer(container_to_kill, 1)
	if err != nil {
		return err
	}

	removeOpts := docker.RemoveContainerOptions{
		ID: container_to_kill,
	}

	err = c.client.RemoveContainer(removeOpts)
	return err
}

func (c *DockerScheduler) StopStartTasks(tasks map[string]demand.Task, ready chan struct{}) error {
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

	go func() {

		// Scale down first to free up resources
		for _, name := range too_many {
			task := tasks[name]
			diff = task.Requested - task.Demand
			log.Printf("Stop %d of task %s", diff, name)
			for i := 0; i < diff; i++ {
				err = c.stopTask(name, &task)
				if err != nil {
					log.Printf("Couldn't stop %s: %v ", name, err)
				}
				task.Requested -= 1
			}
			tasks[name] = task
		}

		// Now we can scale up
		for _, name := range too_few {
			task := tasks[name]
			diff = task.Demand - task.Requested
			log.Printf("Start %d of task %s", diff, name)
			for i := 0; i < diff; i++ {
				err = c.startTask(name, &task)
				if err != nil {
					log.Printf("Couldn't start %s: %v ", name, err)
				}
				task.Requested += 1
			}
			tasks[name] = task
		}

		// Notify the channel when the scaling command has completed
		ready <- struct{}{}
	}()

	return err
}

func (c *DockerScheduler) CountAllTasks(tasks map[string]demand.Task) error {
	// Docker Remote API https://docs.docker.com/reference/api/docker_remote_api_v1.20/
	// get /containers/json
	var err error
	var containers []docker.APIContainers
	containers, err = c.client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return fmt.Errorf("Failed to list containers: %v", err)
	}

	// Reset all the running counts to 0
	for name, t := range tasks {
		t.Running = 0
		tasks[name] = t
	}

	var service_name string
	var present bool

	for i := range containers {
		labels := containers[i].Labels
		service_name, present = labels[f12_map]
		if present {
			// Only update tasks that are already in our task map - don't try to manage anything else
			// log.Printf("Found a container with labels %v", labels)
			t, in_our_tasks := tasks[service_name]
			if in_our_tasks {
				t.Running++
				tasks[service_name] = t
			}
		}
	}

	return err
}
