package docker

import (
	"fmt"
	"testing"

	"github.com/force12io/force12/demand"
)

func TestDockerScheduler(t *testing.T) {
	d := NewScheduler()
	var task demand.Task
	task.Demand = 5
	task.Image = "force12io/force12-demo:latest"

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
