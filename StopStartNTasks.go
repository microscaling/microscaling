package main

import (
	//"fmt"
	"time"
	//"sync"
	"log"
	"strconv"
	"strings"
	//"math/rand"
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
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
func StopStartNTasks(app string, family string, demandcount int, currentcount int, force bool) {
	// Submit a post request to Marathon to match the requested number of the requested app
	// format looks like:
	// PUT http://marathon.force12.io:8080/v2/apps/<app>
	//  Request:
	//  {
	//    "instances": 8
	//  }
	var str string
	var port string
	port = os.Getenv("MARATHON_PORT")
	str = os.Getenv("MARATHON_ADDRESS")
	str = str + port
	if port == "" {
		port = "8080"
	}
	if str == "" {
		str = "http://marathon.force12.io:" + port
	}
	str += "/v2/apps/"
	str += app
	
	if force {
	  str += "?force=true"
	}
	log.Println("Start/stop PUT: " + str)

	var jsonStr string
	jsonStr = "{\"instances\":xxxxxxxxxx}"
	jsonStr = strings.Replace(jsonStr, "xxxxxxxxxx", strconv.Itoa(demandcount), 1)
	log.Println("Start/stop request: " + jsonStr)
	
	//req.Header.Set("X-Custom-Header", "myvalue")
	//req.Header.Set("Content-Type", "application/json")
	var query = []byte(jsonStr)
	req, err1 := http.NewRequest("PUT", str, bytes.NewBuffer(query))

	if err1 != nil {
		log.Println("NewRequest err")
	}

	client := &http.Client{}
	resp1, err1 := client.Do(req)
	if resp1 != nil {
		defer resp1.Body.Close()
	}
	if err1 != nil || resp1 == nil {
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
			// Check for an error in the response that looks like that
      //    "message": "App is locked by one or more deployments. Override with the option '?force=true'. View details at '/v2/deployments/<DEPLOYMENT_ID>'.",
      //    "deployments": [
      //    {
      //      "id": "823714e0-f36e-4401-bcb6-13cf5e05ca04"
      //    }
      //    ]
      var json_prefix string = "App is locked"
	    stringslice := strings.Split(s, json_prefix)

	    if len(stringslice) >= 2 && force == false {
	      // don't force if we have already tried forcing
		    log.Println("App is locked, force it")
		    var sleep time.Duration
		    sleepcount, errenv := strconv.Atoi(os.Getenv("SLEEP_BEFORE_FORCE"))
		    if errenv != nil {
		      sleepcount = 50
		    }
			  sleep = time.Duration(sleepcount) * time.Millisecond
			  time.Sleep(sleep)
		    StopStartNTasks(app, family, demandcount, currentcount, true)
		  }
	  }
	}
	return
}
