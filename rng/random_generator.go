// rng generates a random demand metric for high priority containers
package rng

import (
	"fmt"
	"log"
	"math/rand"

	"bitbucket.org/force12io/force12-scheduler/demand"
)

// TODO! The maximum should not be defined here
const maximum int = 9 // Demand can vary between 0 and maximum
const delta int = 3   // Current value can only go up or down by a maximum of delta

type RandomDemand struct {
	currentP1Demand int
}

// check that we implement the demand interface
var _ demand.Input = (*RandomDemand)(nil)

// NewDemandModel created a new RNG demand model
func NewDemandModel() *RandomDemand {
	return &RandomDemand{
		currentP1Demand: 0,
	}
}

// GetDemand generates the demand, which will be within +/- delta of the current value.
// At present taskType must be "priority1-demand". No other values are supported.
func (rng *RandomDemand) GetDemand(taskType string) (int, error) {
	var newDemand int
	var err error = nil

	switch taskType {
	case "priority1-demand": // TODO! Priority name shouldn't be hard-coded like this

		// Random value between +/- delta is the same as
		// (random value between 0 and 2*delta) - delta
		// noting that if r = rand.Intn(n) then 0 <= r < n
		r := rand.Intn((2 * delta) + 1)
		newDemand = rng.currentP1Demand + r - delta
		if newDemand > maximum {
			newDemand = maximum
		}

		if newDemand < 0 {
			newDemand = 0
		}

		log.Printf("P1 random demand %d", newDemand)
		rng.currentP1Demand = newDemand
	default:
		err = fmt.Errorf("Wrong task type passed to RNG: %s", taskType)
	}

	return newDemand, err
}
