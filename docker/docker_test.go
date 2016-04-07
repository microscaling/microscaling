package docker

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	// "strings"
	"testing"

	// "github.com/microscaling/microscaling/api/apitest"
	"github.com/fsouza/go-dockerclient"
	"github.com/microscaling/microscaling/demand"
)

func TestDockerInitScheduler(t *testing.T) {
	tests := []struct {
		pullImages bool
	}{
		{
			pullImages: true,
		},
		{
			pullImages: false,
		},
	}

	for _, test := range tests {
		d := NewScheduler(test.pullImages, "unix:///var/run/docker.sock")
		log.Printf("Should I pull images? %v", test.pullImages)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Received something %v", r)

			// TODO!! Check that we're receiving what we expect
		}))

		log.Printf("Test server at %s", server.URL)
		d.client, _ = docker.NewClient(server.URL)
		log.Printf("Docker client %v", d.client)

		var task demand.Task
		task.Demand = 5
		task.Image = "microscaling/priority-1:latest"

		d.InitScheduler("anything", &task)
		err := d.startTask("anything", &task)
		if err != nil {
			fmt.Printf("Error %v", err)
		}
	}
}

func TestDockerScheduler(t *testing.T) {
	d := NewScheduler(true, "unix:///var/run/docker.sock")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	d.client, _ = docker.NewClient(server.URL)

	var task demand.Task
	task.Demand = 5
	task.Image = "microscaling/priority-1:latest"

	d.InitScheduler("anything", &task)

	err := d.startTask("anything", &task)
	if err != nil {
		// We don't actually expect these to work locally
		// TODO! Some Docker tests that mock out the Docker client
		fmt.Printf("Error %v", err)
	}

	var tasks map[string]demand.Task
	tasks = make(map[string]demand.Task)
	tasks["anything"] = task
	d.CountAllTasks(tasks)
}
