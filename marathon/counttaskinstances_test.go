package marathon

import (
	"testing"

	"bitbucket.org/force12io/force12-scheduler/demand"
)

func TestCountTasks(t *testing.T) {
	var tasks map[string]demand.Task = make(map[string]demand.Task)

	tasks["anything"] = demand.Task{Demand: 8, Requested: 3}
	tasks["something"] = demand.Task{Demand: 0, Requested: 1}
	m := NewScheduler()

	_ = m.CountAllTasks(tasks)

	for name, task := range tasks {
		if task.Running != task.Requested {
			t.Fatalf("Task %s running is not what was requested or demanded", name)
		}
	}
}
