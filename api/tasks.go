// API between Force12 agent and server
package api

import (
	"encoding/json"
	"log"
)

type TaskDescription struct {
	App         string `json:"app"`
	DemandCount int    `json:"demandCount"`
}

// Get /tasks/ to receive the current task demand for app
func GetTasks(userID string) (td []TaskDescription, err error) {
	body, err := getJsonGet(userID, "/tasks/")
	if err != nil {
		log.Printf("Failed to get /tasks/: %v", err)
		return nil, err
	}

	err = json.Unmarshal(body, &td)
	return td, err
}
