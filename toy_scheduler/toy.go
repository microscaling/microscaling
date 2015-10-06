// Toy scheduler is a mock scheduling output that simply reflects back whatever we tell it
package toy_scheduler

import (
	"fmt"
	"log"

	"bitbucket.org/force12io/force12-scheduler/demand"
	"bitbucket.org/force12io/force12-scheduler/scheduler"
)

type ToyScheduler struct {
}

func NewScheduler() *ToyScheduler {
	toy := ToyScheduler{}
	return &toy
}

// compile-time assert that we implement the right interface
var _ scheduler.Scheduler = (*ToyScheduler)(nil)

func (t *ToyScheduler) InitScheduler(appId string, task *demand.Task) error {
	log.Printf("Toy scheduler initialized task %s with %d initial demand", appId, task.Demand)
	return nil
}

// StopStartNTasks asks the scheduler to bring the number of running tasks up to task.Demand.
func (t *ToyScheduler) StopStartNTasks(appId string, task *demand.Task) error {
	if appId == "force12" {
		return fmt.Errorf("Don't try to scale our own force12 task!")
	}
	task.Requested = task.Demand
	return nil
}

// CountAllTasks for the Toy scheduler simply reflects back what has been requested
func (t *ToyScheduler) CountAllTasks(tasks map[string]demand.Task) error {
	for name, task := range tasks {
		task.Running = task.Requested
		tasks[name] = task
	}
	return nil
}
