package marathon

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type startStopPayload struct {
	Instances int `json:"instances"`
}

////////////////////////////////////////////////////////////////////////////////////////
// StopStartNTasks
//
// This function calls the Marathon API to create the number of app instances (containers)
// we want. For consistency with the ECS API we take the current count, but actually
// we don't use it as Marathon will just work out how many to start and stop based on
// what we tell it we need.
//
//
func StopStartNTasks(app string, family string, demandcount int, currentcount int) {
	// Submit a post request to Marathon to match the requested number of the requested app
	// format looks like:
	// PUT http://marathon.force12.io:8080/v2/apps/<app>
	//  Request:
	//  {
	//    "instances": 8
	//  }
	url := getBaseMarathonUrl() + "/" + app
	log.Println("Start/stop PUT: " + url)

	payload := startStopPayload{
		Instances: demandcount,
	}
	w := &bytes.Buffer{}
	encoder := json.NewEncoder(w)
	err := encoder.Encode(&payload)
	if err != nil {
		log.Printf("Failed to encode json. %v")
		return
	}

	req, err := http.NewRequest("PUT", url, w)

	if err != nil {
		log.Println("NewRequest err %v", err)
	}
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		// handle error
		log.Println("start/stop err %v", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		log.Printf("start/stop read err %v", err)
		return
	}

	s := string(body)
	log.Printf("start/stop json: %s", s)

	return
}
