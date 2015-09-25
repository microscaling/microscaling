package main

import (
	// "log"
	"testing"

	"bitbucket.org/force12io/force12-scheduler/demand"
	"bitbucket.org/force12io/force12-scheduler/rng"
)

func TestDemandUpdate(t *testing.T) {
	var demandchange bool

	tasks = make(map[string]demand.Task)
	tasks["priority1"] = demand.Task{
		FamilyName: "p1family",
		Demand:     const_p1demandstart,
		Requested:  0,
	}

	tasks["priority2"] = demand.Task{
		FamilyName: "p2family",
		Demand:     const_p2demandstart,
		Requested:  0,
	}

	di := rng.NewDemandModel()

	demandchange, _ = update(di)
	if !demandchange {
		// Note this test relies on us not seeding random numbers. Not very nice but OK for our purposes.
		t.Fatalf("Expected demand to have changed but it didn't")
	}
}
