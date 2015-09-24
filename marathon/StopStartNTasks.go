package marathon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type startStopPayload struct {
	Instances int `json:"instances"`
}

// StopStartNTasks calls the Marathon API to create the number of app instances (containers)
// we want. For Marathon we don't need the current count (as we do for ECS API)
// as Marathon will just work out how many to start and stop based on what we tell it we need.
func (m *MarathonScheduler) StopStartNTasks(app string, family string, demandcount int, currentcount int) error {
	// Submit a post request to Marathon to match the requested number of the requested app
	// format looks like:
	// PUT http://marathon.force12.io:8080/v2/apps/<app>
	//  Request:
	//  {
	//    "instances": 8
	//  }
	url := m.baseMarathonUrl + "/" + app
	log.Printf("Start/stop PUT: %s", url)

	payload := startStopPayload{
		Instances: demandcount,
	}
	w := &bytes.Buffer{}
	encoder := json.NewEncoder(w)
	err := encoder.Encode(&payload)
	if err != nil {
		return fmt.Errorf("Failed to encode json. %v", err)
	}

	req, err := http.NewRequest("PUT", url, w)
	if err != nil {
		return fmt.Errorf("Failed to build PUT request err %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return fmt.Errorf("start/stop err %v", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("error response from marathon. %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("start/stop read err %v", err)
	}

	// We do nothing with this body
	s := string(body)
	log.Printf("start/stop json: %s", s)

	return nil
}
