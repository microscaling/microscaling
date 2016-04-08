package main

import (
	"github.com/microscaling/microscaling/api"
	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/scheduler"
)

// handleDemandChange updates to changed demand
func handleDemandChange(td []api.TaskDemand, s scheduler.Scheduler, running *demand.Tasks) (err error) {
	running.Lock()
	defer running.Unlock()
	runningTasks := running.Tasks

	var demandChanged = false
	for _, task := range td {
		name := task.App

		if existingTask, ok := runningTasks[name]; ok {
			if existingTask.Demand != task.DemandCount {
				demandChanged = true
			}
			existingTask.Demand = task.DemandCount
			runningTasks[name] = existingTask
		}
	}

	if demandChanged {
		// Ask the scheduler to make the changes
		err = s.StopStartTasks(running.Tasks)
		if err != nil {
			log.Errorf("Failed to stop / start tasks. %v", err)
		}
	}

	return
}
