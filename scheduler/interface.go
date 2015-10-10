// scheduler defines the interface for scheduling models
package scheduler

import (
	"bitbucket.org/force12io/force12-scheduler/demand"
)

type Scheduler interface {
	// InitScheduler creates and starts the app identified by appId
	InitScheduler(appId string, task *demand.Task) error

	// StopStartNTasks changes the count of containers to match task.Demand
	// returns true if we can immediately call this again, otherwise
	StopStartNTasks(appId string, task *demand.Task, ready chan struct{}) error

	// TODO! It might be more efficient to have a start-stop that can change counts for multiple
	// container types at once

	// CountAllTasks updates task.Running to tell us how many instances of each task are currently running
	CountAllTasks(tasks map[string]demand.Task) error
}
