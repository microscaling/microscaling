package main

import (
	"log"

	"github.com/microscaling/microscaling/api"
	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/scheduler"
)

// handleDemandChange updates to changed demand
func handleDemandChange(td []api.TaskDemand, s scheduler.Scheduler, tasks map[string]demand.Task) (err error) {
	var demandChanged = false
	for _, task := range td {
		name := task.App

		if existingTask, ok := tasks[name]; ok {
			if existingTask.Demand != task.DemandCount {
				demandChanged = true
			}
			existingTask.Demand = task.DemandCount
			tasks[name] = existingTask
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
