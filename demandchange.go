package main

import (
	"log"

	"github.com/force12io/force12/api"
	"github.com/force12io/force12/demand"
	"github.com/force12io/force12/scheduler"
)

// handleDemandChange updates to changed demand
func handleDemandChange(td []api.TaskDemand, s scheduler.Scheduler, tasks map[string]demand.Task) (err error) {
	var demandChanged bool = false
	for _, task := range td {
		name := task.App

		if existing_task, ok := tasks[name]; ok {
			if existing_task.Demand != task.DemandCount {
				demandChanged = true
			}
			existing_task.Demand = task.DemandCount
			tasks[name] = existing_task
		}
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
