// scheduler defines the interface for scheduling models
package scheduler

import (
	"bitbucket.org/force12io/force12-scheduler/demand"
)

type Scheduler interface {
	// InitScheduler creates and starts the app identified by appId
	InitScheduler(appId string) error

	// StopStartNTasks changes the count of containers for the app from currentcount
	// to demandcount
	StopStartNTasks(appId string, family string, demandcount int, currentcount *int) error

	// CountTaskInstances tells us how many instances of this task are currently running / requested
	CountTaskInstances(taskName string, task demand.Task) (int, int, error)
}
