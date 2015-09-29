package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"bitbucket.org/force12io/force12-scheduler/demand"
	"bitbucket.org/force12io/force12-scheduler/scheduler"
)

// TODO! Make this less specific to P1 & P2 model so it returns json structure for arbitrary set of tasks
type sendStatePayload struct {
	CreatedAt          int64 `json:"createdAt"`
	Priority1Requested int   `json:"priority1Requested"`
	Priority1Running   int   `json:"priority1Running"`
	Priority2Running   int   `json:"priority2Running"`
}

func getBaseF12APIUrl() string {
	baseUrl := os.Getenv("F12_API_ADDRESS")
	if baseUrl == "" {
		baseUrl = "https://force12-windtunnel.herokuapp.com"
	}
	return baseUrl
}

var baseF12APIUrl string = getBaseF12APIUrl()

// sendState checks the current state of tasks and sends that state to the f12 API
func SendState(userID string, sched scheduler.Scheduler, tasks map[string]demand.Task) error {
	var err error = nil

	for name, task := range tasks {
		t := task
		t.Running, _, err = sched.CountTaskInstances(name, t)
		if err != nil {
			return fmt.Errorf("Failed to count %s tasks %v", name, err)
		}
		tasks[name] = t
	}

	// Submit a PUT request to the API
	url := baseF12APIUrl + "/metrics/" + userID
	log.Printf("API PUT: %s", url)

	// TODO! Make this less specific to P1 & P2 model
	payload := sendStatePayload{
		CreatedAt:          time.Now().Unix(),
		Priority1Requested: tasks["priority1"].Demand,
		Priority1Running:   tasks["priority1"].Running,
		Priority2Running:   tasks["priority2"].Running,
	}

	w := &bytes.Buffer{}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&payload)
	if err != nil {
		return fmt.Errorf("Failed to encode API json. %v", err)
	}

	req, err := http.NewRequest("PUT", url, w)

	if err != nil {
		return fmt.Errorf("Failed to build API PUT request err %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return fmt.Errorf("API state err %v", err)
	}

	if resp.StatusCode > 204 {
		return fmt.Errorf("error response from API. %s", resp.Status)
	}
	return err
}
