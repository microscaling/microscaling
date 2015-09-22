// Generate a random demand metric for high priority containers
package rng

import (
  "log"
  "math/rand"
)

const maximum int = 9 // Demand can vary between 0 and maximum
const delta int = 3 // Current value can only go up or down by a maximum of delta


func RandomDemand(current_demand int) (int, error) {
  // Random value between +/- delta is the same as 
  // (random value between 0 and 2*delta) - delta
  // noting that if r = rand.Intn(n) then 0 <= r < n 

  r := rand.Intn((2 * delta) + 1)
  demand := current_demand + r - delta
  if demand > maximum {
    demand = maximum
  }

  if demand < 0 {
    demand = 0
  }

  log.Printf("Random demand %d", demand)
  return demand, nil
}