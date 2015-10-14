package toy_scheduler

import (
	"log"
	"testing"

	"bitbucket.org/force12io/force12-scheduler/demand"
)

func TestToyScheduler(t *testing.T) {
	var tasks map[string]demand.Task = make(map[string]demand.Task)

	tasks["anything"] = demand.Task{Demand: 8, Requested: 3}
	m := NewScheduler()

	task := tasks["anything"]
	m.InitScheduler("anything", &task)

	log.Printf("before start/stop: demand %d, requested %d, running %d", task.Demand, task.Requested, task.Running)
	ready := make(chan struct{}, 1)
	err := m.StopStartTasks(tasks, ready)
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	task = tasks["anything"]
	log.Printf("after start/stop: demand %d, requested %d, running %d", task.Demand, task.Requested, task.Running)

	if err != nil {
		t.Fatalf("Error. %v", err)
	} else if task.Requested != task.Demand {
		t.Fatalf("Requested should have been updated")
	}

	err = m.CountAllTasks(tasks)
	for name, task := range tasks {
		if task.Running != task.Requested || task.Running != task.Demand {
			t.Fatalf("Task %s running is not what was requested or demanded", name)
		}
		log.Printf("after counting: demand %d, requested %d, running %d", task.Demand, task.Requested, task.Running)

	}

}
