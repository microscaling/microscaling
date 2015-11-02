package main

import (
	"log"
	"testing"

	"github.com/force12io/force12/demand"
	"github.com/force12io/force12/rng"
	"github.com/force12io/force12/toy_scheduler"
)

func TestDemandUpdate(t *testing.T) {
	var demandchange bool

	tasks = make(map[string]demand.Task)
	tasks["priority1"] = demand.Task{
		FamilyName: "p1family",
		Demand:     3,
		Requested:  0,
	}

	tasks["priority2"] = demand.Task{
		FamilyName: "p2family",
		Demand:     5,
		Requested:  0,
	}

	di := rng.NewDemandModel(4, 10)

	demandchange, _ = update(di, tasks)
	if !demandchange {
		// Note this test relies on us not seeding random numbers. Not very nice but OK for our purposes.
		t.Fatalf("Expected demand to have changed but it didn't")
	}
}

func TestHandleDemandChange(t *testing.T) {
	tasks = make(map[string]demand.Task)
	tasks["priority1"] = demand.Task{
		FamilyName: "p1family",
		Demand:     4,
		Requested:  0,
	}

	tasks["priority2"] = demand.Task{
		FamilyName: "p2family",
		Demand:     3,
		Requested:  0,
	}

	// We might see our own task when we look at Docker, we shouldn't be scaling it!
	tasks["force12"] = demand.Task{
		FamilyName: "force12",
		Demand:     1,
		Requested:  1,
	}

	di := rng.NewDemandModel(3, 9)
	s := toy_scheduler.NewScheduler()

	for i := 0; i < 5; i++ {
		err := handleDemandChange(di, s, tasks)
		if err != nil {
			t.Fatalf("handleDemandChange failed")
		}
		log.Println(tasks)
	}
}
