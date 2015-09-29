package marathon

import (
	"bitbucket.org/force12io/force12-scheduler/demand"
)

func (m *MarathonScheduler) CountAllTasks(tasks map[string]demand.Task) error {
	for name, task := range tasks {
		task.Running = task.Requested
		tasks[name] = task
	}
	return nil
}
