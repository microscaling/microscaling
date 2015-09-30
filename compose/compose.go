// Schedule using docker compose
package compose

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"bitbucket.org/force12io/force12-scheduler/demand"
	"bitbucket.org/force12io/force12-scheduler/scheduler"

	"github.com/fsouza/go-dockerclient"
)

type ComposeScheduler struct {
	client *docker.Client
}

type ComposeContainer struct {
	id    string   `json:"Id"`
	names []string `json:"Names"`
}

func NewScheduler() *ComposeScheduler {
	client, _ := docker.NewClient(os.Getenv("DOCKER_HOST"))

	return &ComposeScheduler{
		client: client,
	}
}

// compile-time assert that we implement the right interface
var _ scheduler.Scheduler = (*ComposeScheduler)(nil)

func (c *ComposeScheduler) InitScheduler(appId string) error {
	// Nothing to do here. yaml file from windtunnel will start one container of each type
	return nil
}

func (c *ComposeScheduler) StopStartNTasks(appId string, task *demand.Task) error {
	// Shell out to Docker compose scale
	// docker-compose scale web=2 worker=3

	var err error

	param := fmt.Sprintf("%s=%d", appId, task.Demand)
	cmd := exec.Command("docker-compose", "scale", param)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("Stderr: %s", stderr.String())
		return err
	}

	task.Requested = task.Demand

	return err
}

func (c *ComposeScheduler) CountAllTasks(tasks map[string]demand.Task) error {

	// Docker Remote API https://docs.docker.com/reference/api/docker_remote_api_v1.20/
	// get /containers/json
	var err error
	var containers []docker.APIContainers
	containers, err = c.client.ListContainers(docker.ListContainersOptions{All: true})
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
		service_name, present = labels["com.docker.compose.service"]
		if present {
			t := tasks[service_name]
			t.Running++
			tasks[service_name] = t
		}
	}

	fmt.Println(tasks)
	return err
}
