package toy_scheduler

import (
	"log"
	"testing"

	"bitbucket.org/force12io/force12-scheduler/demand"
)

func TestToyScheduler(t *testing.T) {
	toy := NewScheduler()
	log.Println(toy)

	var task demand.Task = demand.Task{Demand: 8, Requested: 3}

	running, requested, _ := toy.CountTaskInstances("anything", task)
	if running != 3 {
		t.Fatalf("Wrong running count")
	}

	if requested != 3 {
		t.Fatalf("Wrong requested count")
	}
}

func TestStartStop(t *testing.T) {

	var task demand.Task = demand.Task{Demand: 8, Requested: 3}
	m := NewScheduler()

	log.Println("before start/stop: current, demand", task.Demand, task.Requested)
	err := m.StopStartNTasks("blah", "foobar", task.Demand, &task.Requested)
	log.Println("after start/stop: current, demand", task.Demand, task.Requested)

	if err != nil {
		t.Fatalf("Error. %v", err)
	} else if task.Requested != task.Demand {
		t.Fatalf("Requested should have been updated")
	}
}
