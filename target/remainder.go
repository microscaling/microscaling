package target

// RemainderTarget is what we use as a target when we want to use any remaining capacity for this task
type RemainderTarget struct {
	maxContainers int
}

// NewRemainderTarget creates a new remainder target
func NewRemainderTarget(maxContainers int) *RemainderTarget {
	return &RemainderTarget{
		maxContainers: maxContainers,
	}
}

// Meeting returns true if the target is currently met
func (t *RemainderTarget) Meeting(current int) bool {
	// log.Debugf("[remainder] meeting: always false")
	return false
}

// Exceeding returns true if the target is currently exceeded
func (t *RemainderTarget) Exceeding(current int) bool {
	// log.Debugf("[remainder] exceeding: always false")
	return false
}

// Delta returns the nuumber of additional containers we should add. For a remainder target we always try to get as many as we can.
func (t *RemainderTarget) Delta(current int) (delta int) {
	delta = t.maxContainers
	return
}
