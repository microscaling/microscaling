package docker

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
		log.Infof("Should I pull images? %v", test.pullImages)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Infof("Received something %v", r)

			// TODO!! Check that we're receiving what we expect
		}))

		log.Infof("Test server at %s", server.URL)
		d.client, _ = docker.NewClient(server.URL)
		log.Debugf("Docker client %v", d.client)

		var task demand.Task
		task.Demand = 5
		task.Image = "microscaling/priority-1:latest"

		d.InitScheduler("anything", &task)
		d.startTask("anything", &task)
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

	d.startTask("anything", &task)
	// TODO! Some Docker tests that mock out the Docker client

	var tasks demand.Tasks
	tasks.Tasks = make(map[string]demand.Task)
	tasks.Tasks["anything"] = task
	d.CountAllTasks(&tasks)
}
