package toy

import (
	"testing"

	"github.com/microscaling/microscaling/demand"
)

func TestToyScheduler(t *testing.T) {
	var tasks demand.Tasks
	tasks.Tasks = make(map[string]demand.Task)

	tasks.Tasks["anything"] = demand.Task{Demand: 8, Requested: 3}
	m := NewScheduler()

	task := tasks.Tasks["anything"]
	m.InitScheduler("anything", &task)

	log.Debugf("before start/stop: demand %d, requested %d, running %d", task.Demand, task.Requested, task.Running)
	err := m.StopStartTasks(tasks.Tasks)
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	task = tasks.Tasks["anything"]
	log.Debugf("after start/stop: demand %d, requested %d, running %d", task.Demand, task.Requested, task.Running)

	if err != nil {
		t.Fatalf("Error. %v", err)
	} else if task.Requested != task.Demand {
		t.Fatalf("Requested should have been updated")
	}

	err = m.CountAllTasks(&tasks)
	for name, task := range tasks.Tasks {
		if task.Running != task.Requested || task.Running != task.Demand {
			t.Fatalf("Task %s running is not what was requested or demanded", name)
		}
		log.Debugf("after counting: demand %d, requested %d, running %d", task.Demand, task.Requested, task.Running)

	}

}
