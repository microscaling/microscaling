package rng

import (
	"math"
	"testing"
)

func TestRandomDemand(t *testing.T) {

	rng := NewRandomDemandGenerator()

	for i := 0; i < 20; i++ {
		old_demand := rng.currentDemand
		demand, _ := rng.GetDemand("anything")

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
}
