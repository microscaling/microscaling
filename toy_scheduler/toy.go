// Toy scheduler is a mock scheduling output that simply reflects back whatever we tell it
package toy_scheduler

import (
	"log"

	"bitbucket.org/force12io/force12-scheduler/demand"
	"bitbucket.org/force12io/force12-scheduler/scheduler"
)

type ToyScheduler struct {
	tasks map[string]demand.Task
}

func NewScheduler() *ToyScheduler {
	toy := ToyScheduler{}
	toy.tasks = make(map[string]demand.Task)
	toy.tasks["priority1"] = demand.Task{
		FamilyName: "p1-family",
		Demand:     5,
		Requested:  0,
	}

	toy.tasks["priority2"] = demand.Task{
		FamilyName: "p2-family",
		Demand:     4,
		Requested:  0,
	}
	return &toy
}

// compile-time assert that we implement the right interface
var _ scheduler.Scheduler = (*ToyScheduler)(nil)

func (t *ToyScheduler) InitScheduler(appId string) error {
	log.Printf("Toy scheduler initialized task %s", appId)
	return nil
}

// StopStartNTasks asks the scheduler to bring the number of running tasks up to demandcount.
func (t *ToyScheduler) StopStartNTasks(appId string, family string, demandcount int, requestedcount *int) error {
	*requestedcount = demandcount
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
