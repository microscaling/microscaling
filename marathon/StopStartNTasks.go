package marathon

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

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
	var str string
	str = os.Getenv("MARATHON_ADDRESS")
	if str == "" {
		str = "http://marathon.force12.io:8080"
	}
	str += "/v2/apps/"
	str += app
	log.Println("Start/stop PUT: " + str)

	var jsonStr string
	jsonStr = "{\"instances\":xxxxxxxxxx}"
	jsonStr = strings.Replace(jsonStr, "xxxxxxxxxx", strconv.Itoa(demandcount), 1)
	log.Println("Start/stop request: " + jsonStr)

	var query = []byte(jsonStr)
	req, err1 := http.NewRequest("PUT", str, bytes.NewBuffer(query))

	if err1 != nil {
		log.Println("NewRequest err")
	}
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp1, err1 := client.Do(req)

	defer resp1.Body.Close()
	if err1 != nil {
		// handle error
		log.Println("start/stop err")
	} else {
		body, err0 := ioutil.ReadAll(resp1.Body)
		if err0 != nil {
			// handle error
			log.Println("start/stop read err")
		} else {
			s := string(body[:])
			log.Println("start/stop json: " + s)
		}
	}
	return
}
