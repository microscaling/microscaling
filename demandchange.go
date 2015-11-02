package main

import (
	"log"

	"github.com/force12io/force12/demand"
	"github.com/force12io/force12/scheduler"
)

// handleDemandChange checks the new demand
func handleDemandChange(input demand.Input, s scheduler.Scheduler, tasks map[string]demand.Task) error {
	var err error = nil
	var demandChanged bool

	demandChanged, err = update(input, tasks)
	if err != nil {
		log.Printf("Failed to get new demand. %v", err)
		return err
	}

	if demandChanged {
		// Ask the scheduler to make the changes
		err = s.StopStartTasks(tasks)
		if err != nil {
			log.Printf("Failed to stop / start tasks. %v", err)
		}
	}

	return err
}

// update checks for changes in demand, returning true if demand changed
func update(input demand.Input, ts map[string]demand.Task) (bool, error) {
	var err error = nil
	var demandchange bool = false

	for name, task := range ts {
		oldDemand := task.Demand
		task.Demand, err = input.GetDemand(name)
		if err != nil {
			log.Printf("Failed to get new demand for task %s. %v", name, err)
			return demandchange, err
		}

		log.Printf("Current demand: task %s - %d", name, task.Demand)

		if task.Demand != oldDemand {
			demandchange = true
		}

		ts[name] = task
	}
	return demandchange, err
}
