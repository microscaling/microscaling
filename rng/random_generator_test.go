package rng

import (
	"log"
	"math"
	"testing"
)

func TestRandomDemand(t *testing.T) {

	maximum := 9
	delta := 3
	rng := NewDemandModel(delta, maximum)

	for i := 0; i < 20; i++ {
		old_demand := rng.currentP1Demand
		demand1, _ := rng.GetDemand("priority1")
		demand2, _ := rng.GetDemand("priority2")
		log.Printf("Demand changed from %d to %d", old_demand, demand1)

		if demand1 > maximum || demand2 > maximum {
			t.Fatalf("Random value exceeds maximum")
		}

		if demand1 < 0 || demand2 < 0 {
			t.Fatalf("Random value below 0")
		}

		if math.Abs(float64(demand1)-float64(old_demand)) > float64(delta) {
			t.Fatalf("Random value varied more than the delta")
		}
	}

	// Right now you should only pass in priority1 or priority2
	_, err := rng.GetDemand("something")
	if err == nil {
		t.Fatalf("Failed to barf on the wrong task type")
	}
}
