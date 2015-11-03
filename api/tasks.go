// API between Force12 agent and server
package api

import (
	"encoding/json"
	"log"

	"github.com/force12io/force12/demand"
)

type TaskDescription struct {
	App         string `json:"app"`
	DemandCount int    `json:"demandCount"`
}

func tasksFromResponse(b []byte, tasks map[string]demand.Task) (err error) {
	var t []TaskDescription

	err = json.Unmarshal(b, &t)

	for _, task := range t {
		name := task.App

		if existing_task, ok := tasks[name]; ok {
			existing_task.Demand = task.DemandCount
			tasks[name] = existing_task
		}
	}

	return
}

// Get /tasks/ to receive the current task demand for app
func GetTasks(userID string, tasks map[string]demand.Task) (err error) {
	body, err := getJsonGet(userID, "/tasks/")
	if err != nil {
		log.Printf("Failed to get /tasks/: %v", err)
		return err
	}

	err = tasksFromResponse(body, tasks)
	return err
}
