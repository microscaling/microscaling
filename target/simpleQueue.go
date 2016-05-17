package target

// SimpleQueueLengthTarget aims to keep the number of items in a queue under a certain length
// This version tends to oscillate too much, so we prefer the version that uses a PD controller (queue.go)
type SimpleQueueLengthTarget struct {
	length    int
	minLength int
}

// NewSimpleQueueLengthTarget creates a new target based on simply adding or removing one container if we're above or below target
func NewSimpleQueueLengthTarget(length int) *SimpleQueueLengthTarget {

	return &SimpleQueueLengthTarget{
		length:    length,
		minLength: int(float64(length) * queueLengthExceedingPercent),
	}
}

// Meeting returns true if the target is currently met
func (t *SimpleQueueLengthTarget) Meeting(current int) bool {
	meeting := (current <= t.length)
	if !meeting {
		log.Debugf("[sql] not meeting: current %d target %d", current, t.length)
	}
	return meeting
}

// Exceeding returns true if the target is currently exceeded
func (t *SimpleQueueLengthTarget) Exceeding(current int) bool {
	exceeding := (current <= t.minLength)
	if exceeding {
		log.Debugf("[sql] exceeding: current %d target %d", current, t.length)
	}
	return exceeding
}

// Delta returns the nuumber of additional containers we should add (remove if negative) to try to attain the target
func (t *SimpleQueueLengthTarget) Delta(currentLength int) (delta int) {

	// Simply increment by one if we're over the target, and decrement if we're under
	if currentLength > t.length {
		delta = 1
	} else if currentLength < t.minLength {
		delta = -1
	} else {
		delta = 0
	}

	log.Debugf("[sql] delta %d", delta)
	return
}
