package rng

import (
	"math"
	"testing"
)

func TestRandomDemand(t *testing.T) {

	rng := NewDemandModel()

	for i := 0; i < 20; i++ {
		old_demand := rng.currentP1Demand
		demand, _ := rng.GetDemand("priority1-demand")

		if demand > maximum {
			t.Fatalf("Random value exceeds maximum")
		}

		if demand < 0 {
			t.Fatalf("Random value below 0")
		}

		if math.Abs(float64(demand)-float64(old_demand)) > float64(delta) {
			t.Fatalf("Random value varied more than the delta")
		}
	}

	// Right now you should only pass in priority1-demand
	_, err := rng.GetDemand("something")
	if err == nil {
		t.Fatalf("Failed to barf on the wrong task type")
	}
}
