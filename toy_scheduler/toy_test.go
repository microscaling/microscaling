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
