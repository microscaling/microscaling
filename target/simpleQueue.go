package target

// Target of keeping the number of items in a queue under a certain length
type SimpleQueueLengthTarget struct {
	length    int
	minLength int
}

// const queueLengthExceedingPercent float64 = 0.7

func NewSimpleQueueLengthTarget(length int) Target {
	// TODO!! Better ways to calculate these heuristics and/or pass them in

	return &SimpleQueueLengthTarget{
		length:    length,
		minLength: int(float64(length) * queueLengthExceedingPercent),
	}
}

func (t *SimpleQueueLengthTarget) Meeting(current int) bool {
	meeting := (current <= t.length)
	if !meeting {
		log.Debugf("[sql] not meeting: current %d target %d", current, t.length)
	}
	return meeting
}

func (t *SimpleQueueLengthTarget) Exceeding(current int) bool {
	exceeding := (current <= t.minLength)
	if exceeding {
		log.Debugf("[sql] exceeding: current %d target %d", current, t.length)
	}
	return exceeding
}

// Number of additional containers
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
