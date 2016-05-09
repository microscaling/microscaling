package demand

import (
	"reflect"

	"github.com/microscaling/microscaling/target"
)

// IsRemainder returns true if this task uses up unused resources
func (t *Task) IsRemainder() bool {
	remainderType := reflect.TypeOf(&target.RemainderTarget{})
	ruleType := reflect.TypeOf(t.Target)
	return ruleType == remainderType
}

// ScaleUpCount tells us how many containers to scale up by
// Call this after IdealContainers has been updated
func (t *Task) ScaleUpCount() (delta int) {
	if t.Target.Meeting(t.Metric.Current()) {
		delta = 0
	} else {
		delta = t.IdealContainers - t.Requested
		log.Debugf("Meeting -> delta %d", delta)
	}

	// Only looking for scaling up
	if delta < 0 {
		delta = 0
	}

	// Make sure we do always have at least the minimum
	if t.Requested+delta < t.MinContainers {
		delta = t.MinContainers - t.Requested
		log.Debugf("Need minimum -> delta %d", delta)
	}

	// But make sure this won't exceed the maximum
	if t.Requested+delta > t.MaxContainers {
		delta = t.MaxContainers - t.Requested
		log.Debugf("Can't exceed max -> delta %d", delta)
	}

	if delta > t.MaxDelta {
		delta = t.MaxDelta
	}

	log.Debugf("  [scaleup] %s delta %d", t.Name, delta)
	return
}

// ScaleDownCount tells us how many we should scale down by
// Call this after IdealContainers has been updated
func (t *Task) ScaleDownCount() (delta int) {

	if t.Target.Exceeding(t.Metric.Current()) {
		delta = t.IdealContainers - t.Requested
		log.Debugf("Exceeding -> delta %d", delta)
	} else {
		delta = 0
	}

	// Only looking for scaling down
	if delta > 0 {
		delta = 0
	}

	// Make sure we do always have at least the minimum
	if t.Requested+delta < t.MinContainers {
		delta = t.MinContainers - t.Requested
		log.Debugf("Need minimum -> delta %d", delta)
	}

	// Make sure this won't exceed the maximum
	if t.Requested+delta > t.MaxContainers {
		delta = t.MaxContainers - t.Requested
		log.Debugf("Can't exceed max -> delta %d", delta)
	}

	if delta < -t.MaxDelta {
		delta = -t.MaxDelta
	}

	log.Debugf("  [scaledown] %s delta %d", t.Name, delta)
	return
}

// CanScaleDown returns the number we could scale down by
func (t *Task) CanScaleDown() int {
	if !t.IsScalable {
		return 0
	}

	if t.Requested <= t.MinContainers {
		return 0
	}
	return t.Requested - t.MinContainers
}
