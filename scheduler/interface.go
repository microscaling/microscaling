package scheduler

type Scheduler interface {
	// InitScheduler creates and starts the app identified by appId
	InitScheduler(appId string) error

	// GetContainerCount gets the current count of containers for the app identified by
	// key.
	GetContainerCount(key string) (int, error)

	// StopStartNTasks changes the count of containers for the app from currentcount
	// to demandcount
	StopStartNTasks(appId string, family string, demandcount int, currentcount int) error

	// CountAllTasks is very useful
	CountAllTasks() (int, int, error)
}
