// rng generates a random demand metric for high priority containers
package rng

import (
	"fmt"
	// "log"
	"math/rand"
	"time"

	"github.com/force12io/force12/demand"
)

type RandomDemand struct {
	currentP1Demand int
	delta           int
	maximum         int
}

// check that we implement the demand interface
var _ demand.Input = (*RandomDemand)(nil)

// NewDemandModel creates a new RNG demand model.
// delta - current demand for priority 1 can only go up or down by a maximum of delat
// maximum - max total number of containers
// interval - the number of milliseconds that has to pass before demand is allowed to change
func NewDemandModel(delta int, maximum int) *RandomDemand {

	rand.Seed(int64(time.Now().Nanosecond()))

	return &RandomDemand{
		currentP1Demand: 0,
		delta:           delta,
		maximum:         maximum,
	}
}

// GetDemand generates the demand, which will be within +/- delta of the current value, up to the maximum.
func (rng *RandomDemand) GetDemand(taskType string) (int, error) {
	var newDemand int
	var err error = nil

	switch taskType {
	case "priority1": // TODO! Priority name shouldn't be hard-coded like this
		// Random value between +/- delta is the same as
		// (random value between 0 and 2*delta) - delta
		// noting that if r = rand.Intn(n) then 0 <= r < n
		r := rand.Intn((2 * rng.delta) + 1)
		newDemand = rng.currentP1Demand + r - rng.delta
		if newDemand > rng.maximum {
			newDemand = rng.maximum
		}

		if newDemand < 0 {
			newDemand = 0
		}

		rng.currentP1Demand = newDemand
	case "priority2":
		// priority2 gets whatever is left over
		newDemand = rng.maximum - rng.currentP1Demand
	default:
		err = fmt.Errorf("Wrong task type passed to RNG: %s", taskType)
	}

	return newDemand, err
}
