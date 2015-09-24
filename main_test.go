package main

import (
	"log"
	"testing"

	"bitbucket.org/force12io/force12-scheduler/rng"
)

func TestDemandUpdate(t *testing.T) {
	currentdemand := Demand{}

	di := rng.NewDemandModel()

	log.Println("Demand before: ", currentdemand)
	demandchange := currentdemand.update(di)
	if !demandchange {
		// Note this test relies on us not seeding random numbers. Not very nice but OK for our purposes.
		t.Fatalf("Expected demand to have changed but it didn't")
	}

	log.Println("Demand after: ", currentdemand)
}
