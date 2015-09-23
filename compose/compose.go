// Schedule using docker compose
package compose

import (
	"bitbucket.org/force12io/force12-scheduler/scheduler"
)

type ComposeScheduler struct {
}

func NewScheduler() *ComposeScheduler {
	return &ComposeScheduler{}
}

// compile-time assert that we implement the right interface
var _ scheduler.Scheduler = (*ComposeScheduler)(nil)

func (c *ComposeScheduler) InitScheduler(appId string) error {
	return nil
}

func (c *ComposeScheduler) GetContainerCount(key string) (int, error) {
	return 0, nil
}

func (c *ComposeScheduler) StopStartNTasks(appId string, family string, demandcount int, currentcount int) error {
	return nil
}

func (c *ComposeScheduler) CountAllTasks() (int, int, error) {
	return 0, 0, nil
}
