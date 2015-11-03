package main

import (
	"log"

	"github.com/force12io/force12/demand"
	"github.com/force12io/force12/scheduler"
)

// handleDemandChange checks the new demand
func handleDemandChange(input demand.Input, s scheduler.Scheduler, tasks map[string]demand.Task) (demandChanged bool, err error) {
	demandChanged, err = input.Update(tasks)
	if err != nil {
		log.Printf("Failed to get new demand. %v", err)
		return
	}

	if demandChanged {
		// Ask the scheduler to make the changes
		err = s.StopStartTasks(tasks)
		if err != nil {
			log.Printf("Failed to stop / start tasks. %v", err)
		}
	}

	return
}
