package rng

import (
  "testing"
  "math"
)

func TestRandomDemand(t *testing.T) {
  c := 0
  for i := 0; i < 20; i++ {
    c = randomDemandCall(t, c)
  }
}

func randomDemandCall(t *testing.T, current_demand int) int {
  value, _ := RandomDemand(current_demand)

  if value > maximum {
    t.Fatalf("Random value exceeds maximum")
  }

  if value < 0 {
    t.Fatalf("Random value below 0")
  }

  if math.Abs(float64(value) - float64(current_demand)) > float64(delta) {
    t.Fatalf("Random value varied more than the delta")
  }

  return value
}