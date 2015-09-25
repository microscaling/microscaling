// scheduler defines the interface for scheduling models
package scheduler

type Scheduler interface {
	// InitScheduler creates and starts the app identified by appId
	InitScheduler(appId string) error

	// StopStartNTasks changes the count of containers for the app from currentcount
	// to demandcount
	StopStartNTasks(appId string, family string, demandcount int, currentcount int) error

	// CountAllTasks is very useful
	CountAllTasks() (int, int, error)
}
