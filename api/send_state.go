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
)

// TODO! Make this less specific to P1 & P2 model so it returns json structure for arbitrary set of tasks
type sendStatePayload struct {
	CreatedAt          int64 `json:"createdAt"`
	Priority1Requested int   `json:"priority1Requested"`
	Priority1Running   int   `json:"priority1Running"`
	Priority2Running   int   `json:"priority2Running"`
	MaxContainers      int   `json:"maxContainers"`
}

func getBaseF12APIUrl() string {
	baseUrl := os.Getenv("F12_METRICS_API_ADDRESS")
	if baseUrl == "" {
		baseUrl = "http://app.force12.io"
	}

	log.Printf("Sending results to %s", baseUrl)
	return baseUrl
}

var baseF12APIUrl string = getBaseF12APIUrl()

// sendState sends the current state of tasks to the f12 API
func SendState(userID string, tasks map[string]demand.Task, maxContainers int) error {
	var err error = nil

	// Submit a PUT request to the API
	url := baseF12APIUrl + "/metrics/" + userID
	log.Printf("API PUT: %s | p1 %d running, p2 %d running | p1 demand %d", url,
		tasks["priority1"].Running, tasks["priority2"].Running, tasks["priority1"].Demand)

	// TODO! Make this less specific to P1 & P2 model
	payload := sendStatePayload{
		CreatedAt:          time.Now().Unix(),
		Priority1Requested: tasks["priority1"].Demand,
		Priority1Running:   tasks["priority1"].Running,
		Priority2Running:   tasks["priority2"].Running,
		MaxContainers:      maxContainers,
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
	issuedAt := time.Now()
	resp, err := http.DefaultClient.Do(req)
	apiDuration := time.Since(issuedAt)
	log.Printf("API put took %v", apiDuration)

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
