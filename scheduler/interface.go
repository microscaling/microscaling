// scheduler defines the interface for scheduling models
package scheduler

import (
	"bitbucket.org/force12io/force12-scheduler/demand"
)

type Scheduler interface {
	// InitScheduler creates and starts the app identified by appId
	InitScheduler(appId string, task *demand.Task) error

	// StopStartTasks changes the count of containers to match task.Demand
	// returns true if we can immediately call this again, false if we are waiting for a scale command to complete
	StopStartTasks(tasks map[string]demand.Task, ready chan struct{}) (bool, error)

	// CountAllTasks updates task.Running to tell us how many instances of each task are currently running
	CountAllTasks(tasks map[string]demand.Task) error
}
