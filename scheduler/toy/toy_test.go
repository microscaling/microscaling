package toy

import (
	"testing"

	"github.com/microscaling/microscaling/demand"
)

func TestToyScheduler(t *testing.T) {
	var tasks = demand.Tasks{}
	tasks.Tasks = make([]*demand.Task, 1)

	task := demand.Task{Name: "anything", Demand: 8, Requested: 3}
	tasks.Tasks[0] = &task
	m := NewScheduler()

	m.InitScheduler(&task)

	log.Debugf("before start/stop: demand %d, requested %d, running %d", task.Demand, task.Requested, task.Running)
	err := m.StopStartTasks(&tasks)
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	log.Debugf("after start/stop: demand %d, requested %d, running %d", task.Demand, task.Requested, task.Running)

	if err != nil {
		t.Fatalf("Error. %v", err)
	} else if task.Requested != task.Demand {
		t.Fatalf("Requested should have been updated")
	}

	err = m.CountAllTasks(&tasks)
	for _, task := range tasks.Tasks {
		if task.Running != task.Requested || task.Running != task.Demand {
			t.Fatalf("Task %s running is not what was requested or demanded", task.Name)
		}
		log.Debugf("after counting: demand %d, requested %d, running %d", task.Demand, task.Requested, task.Running)

	}

}
