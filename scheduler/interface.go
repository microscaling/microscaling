// scheduler defines the interface for scheduling models
package scheduler

import (
	"github.com/force12io/force12/demand"
)

type Scheduler interface {
	// InitScheduler creates and starts the app identified by appId
	InitScheduler(appId string, task *demand.Task) error

	// StopStartTasks changes the count of containers to match task.Demand
	// Must notify the ready channel when the tasks are complete
	StopStartTasks(tasks map[string]demand.Task, ready chan struct{}) error

	// CountAllTasks updates task.Running to tell us how many instances of each task are currently running
	CountAllTasks(tasks map[string]demand.Task) error
}
