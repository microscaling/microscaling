package marathon

import (
	"bitbucket.org/force12io/force12-scheduler/demand"
)

// CountTaskInstances for Marathon simply reflects back the number of tasks of this type we have requested
func (m *MarathonScheduler) CountTaskInstances(taskName string, task demand.Task) (int, int, error) {
	return task.Requested, task.Requested, nil
}
