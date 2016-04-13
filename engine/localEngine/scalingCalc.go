package localEngine

import (
	"github.com/microscaling/microscaling/demand"
)

func ScalingCalculation(tasks *demand.Tasks) (demandChanged bool) {
	delta := 0
	demandChanged = false

	// Work out the ideal scale for all the services
	for _, t := range tasks.Tasks {
		t.IdealContainers = t.Running + t.Target.Delta(t.Metric.Current())
	}

	// Look for services we could scale down, in reverse priority order
	tasks.PrioritySort(true)
	for _, t := range tasks.Tasks {
		if !t.IsScalable || t.Requested == t.MinContainers {
			// Can't scale this service down
			continue
		}

		if t.Running != t.Requested {
			// There's a scale operation in progress
			log.Debugf("  [scale] %s already scaling: running %d, requested %d", t.Name, t.Running, t.Requested)
			continue
		}

		// For scaling down, delta should be negative
		delta = t.ScaleDownCount()
		if delta < 0 {
			t.Demand = t.Running + delta
			demandChanged = true
			log.Debugf("  [scale] scaling %s down by %d", t.Name, delta)
		}
	}

	// Now look for tasks we need to scale up
	tasks.PrioritySort(false)
	for p, t := range tasks.Tasks {
		if !t.IsScalable {
			continue
		}

		if t.Running != t.Requested {
			// There's a scale operation in progress
			log.Debugf("  [scale] %s already scaling: running %d, requested %d", t.Name, t.Running, t.Requested)
			continue
		}

		delta = t.ScaleUpCount()
		if delta <= 0 {
			continue
		}

		log.Debugf("  [scale] scale up %s by %d", t.Name, delta)
		available := tasks.CheckCapacity()

		if available < delta {
			// If this is a task that fills the remainder, there's no need to exceed capacity
			if t.IsRemainder() {
				delta = available
			} else {
				log.Debugf("  [scale] looking for %d additional capacity by scaling down:", delta-available)
				index := len(tasks.Tasks)
				freedCapacity := available
				for index > p+1 && freedCapacity < delta {
					// Kill off lower priority services if we need to
					index -= 1
					lowerPriorityService := tasks.Tasks[index]
					if lowerPriorityService.Priority > t.Priority {
						scaleDownBy := lowerPriorityService.CanScaleDown()
						if scaleDownBy > 0 {
							if scaleDownBy > (delta - freedCapacity) {
								scaleDownBy = delta - freedCapacity
							}

							lowerPriorityService.Demand = lowerPriorityService.Running - scaleDownBy
							demandChanged = true
							log.Debugf("  [scale] Service %s priority %d scaling down %d", lowerPriorityService.Name, lowerPriorityService.Priority, -scaleDownBy)
							freedCapacity = freedCapacity + scaleDownBy
						}
					}
				}

				// We might still not have enough capacity and we haven't waited for scale down to complete, so just scale up what's available now
				delta = available
				log.Debugf("  [scale] Can only scale %s by %d", t.Name, delta)
			}
		}

		if delta > 0 {
			demandChanged = true
			if t.Demand >= t.MaxContainers {
				log.Error("*** Not enough capacity for %s ***", t.Name)
				t.Demand = t.MaxContainers
			} else {
				log.Debugf("  [scale] Service %s scaling up %d", t.Name, delta)
				t.Demand = t.Running + delta
			}
		}
	}
	return demandChanged
}
