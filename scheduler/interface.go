// scheduler defines the interface for scheduling models
package scheduler

import (
	"bitbucket.org/force12io/force12-scheduler/demand"
)

type Scheduler interface {
	// InitScheduler creates and starts the app identified by appId
	InitScheduler(appId string, task *demand.Task) error

	// StopStartNTasks changes the count of containers to match task.Demand
	StopStartNTasks(appId string, task *demand.Task) error

	// CountAllTasks updates task.Running to tell us how many instances of each task are currently running
	CountAllTasks(tasks map[string]demand.Task) error
}
