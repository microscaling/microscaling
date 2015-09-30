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
	// client, _ := docker.NewClient("unix:///var/run/docker.sock")
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

func (c *ComposeScheduler) StopStartNTasks(appId string, family string, demandcount int, currentcount *int) error {
	// Shell out to Docker compose scale
	// docker-compose scale web=2 worker=3

	var err error

	param := fmt.Sprintf("%s=%d", appId, demandcount)
	cmd := exec.Command("docker-compose", "scale", param)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("Stderr: %s", stderr.String())
		return err
	}

	*currentcount = demandcount

	return err
}

func (c *ComposeScheduler) CountAllTasks(tasks map[string]demand.Task) error {

	// Docker Remote API https://docs.docker.com/reference/api/docker_remote_api_v1.20/
	// get /containers/json
	var err error

	containers, err := c.client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		return fmt.Errorf("Failed to list containers: %v", err)
	}
	fmt.Println(containers)

	// TODO!! Fill in the running counts for each of the tasks we know about

	return err
}
