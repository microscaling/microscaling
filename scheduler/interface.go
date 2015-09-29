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
	// TODO! This should probably task a *demand.Task rather than individual counts
	StopStartNTasks(appId string, family string, demandcount int, currentcount *int) error

	// CountAllTasks tells us how many instances of each task are currently running
	CountAllTasks(tasks map[string]demand.Task) error
}
