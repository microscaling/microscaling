package target

type RemainderTarget struct {
	maxContainers int
}

func NewRemainderTarget(maxContainers int) Target {
	return &RemainderTarget{
		maxContainers: maxContainers,
	}
}

func (t *RemainderTarget) Meeting(current int) bool {
	// log.Debugf("[remainder] meeting: always false")
	return false
}

func (t *RemainderTarget) Exceeding(current int) bool {
	// log.Debugf("[remainder] exceeding: always false")
	return false
}

// Number of additional containers - we always try to get as many as we can
func (t *RemainderTarget) Delta(current int) (delta int) {
	delta = t.maxContainers
	return
}
