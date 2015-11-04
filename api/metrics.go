// API between Force12 agent and server
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/force12io/force12/demand"
)

type metricsPayload struct {
	CreatedAt int64        `json:"createdAt"`
	Tasks     []appMetrics `json:"tasks"`
}

type appMetrics struct {
	App          string `json:"app"`
	RunningCount int    `json:"runningCount"`
	PendingCount int    `json: "pendingCount"`
}

// sendMetrics sends the current state of tasks to the F12 API
func SendMetrics(userID string, tasks map[string]demand.Task) error {
	var err error = nil
	var index int = 0

	url := baseF12APIUrl + "/metrics/" + userID

	payload := metricsPayload{
		CreatedAt: time.Now().Unix(),
		Tasks:     make([]appMetrics, len(tasks)),
	}

	for name, task := range tasks {
		payload.Tasks[index] = appMetrics{App: name, RunningCount: task.Running, PendingCount: task.Requested}
		index++
	}

	// Submit a PUT request to the API
	w, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to encode API json. %v", err)
	}

	b := bytes.NewBuffer(w)
	req, err := http.NewRequest("PUT", url, b)
	if err != nil {
		return fmt.Errorf("Failed to build API PUT request err %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// issuedAt := time.Now()
	resp, err := httpClient.Do(req)
	// apiDuration := time.Since(issuedAt)
	// log.Printf("API put took %v", apiDuration)

	if err != nil {
		if err.Error() == "EOF" {
			// See http://stackoverflow.com/questions/17714494/golang-http-request-results-in-eof-errors-when-making-multiple-requests-successi
			// We will silently ignore this EOF issue for now
			log.Printf("Ignoring EOF")
			return nil
		}
		return fmt.Errorf("API send state error %v", err)
	}

	if resp == nil || resp.Body == nil {
		log.Printf("Http response is unexpectedly nil")
	}
	defer resp.Body.Close()

	if resp.StatusCode > 204 {
		return fmt.Errorf("error response from API. %s", resp.Status)
	}
	return err
}
