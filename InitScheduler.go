package main

import (
	//"fmt"
	//"time"
	//"sync"
	"log"
	"strings"
	//"strconv"
	//"math/rand"
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
)

////////////////////////////////////////////////////////////////////////////////////////
// InitScheduler
//
// app - string identifier of a container
// Init the Marathon scheduler. For us that means checking whether the supplied app exists on the cluster
// and starting it if not
//
func InitScheduler(app string) {
	// Ask Marathon which apps are running
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
	str += "/v2/apps"

	log.Println("req: " + str)
	resp, err := http.Get(str)
	if resp != nil {
		defer resp.Body.Close()
	}
	log.Println("resp got ")
	if err != nil || resp == nil {
		// handle error
		log.Println("INVALID URL: " + str)
	} else {
		body, err0 := ioutil.ReadAll(resp.Body)
		if err0 != nil {
			// handle error
		} else {
			//Now find out whether both our supplied app is running. If note we need to create one
			//The returned data should be json of the following form in accordance with the Marathon APIs
			// the id field shoudl match the supplied app parameter
			//{
			//"apps": [
			//  {
			//      "id": "/priority1",
			//      "cmd": "/tmp/sleep.sh",
			//      "args": null,
			//      "user": null,
			//      "env": {},
			//      "instances": 2,
			//      "cpus": 0.08,
			//      "mem": 70,
			//      "disk": 0,
			//      "executor": "",
			//      "constraints": [],
			//      "uris": [],
			//      "storeUrls": [],
			//      "ports": [
			//          10000
			//      ],
			//      "requirePorts": false,
			//      "backoffSeconds": 1,
			//      "backoffFactor": 1.15,
			//      "maxLaunchDelaySeconds": 3600,
			//      "container": {
			//          "type": "DOCKER",
			//          "volumes": [],
			//          "docker": {
			//              "image": "quay.io/rossf7/force12-sleeper:latest",
			//              "network": "BRIDGE",
			//              "privileged": false,
			//              "parameters": [],
			//              "forcePullImage": false
			//          }
			//      },
			//      "healthChecks": [],
			//      "dependencies": [],
			//      "upgradeStrategy": {
			//          "minimumHealthCapacity": 1,
			//          "maximumOverCapacity": 1
			//      },
			//      "labels": {},
			//      "acceptedResourceRoles": null,
			//      "version": "2015-07-29T18:03:35.688Z",
			//      "tasksStaged": 0,
			//      "tasksRunning": 2,
			//      "tasksHealthy": 0,
			//      "tasksUnhealthy": 0,
			//      "deployments": []
			//  }
			//]
			//}
			//
			// Rather than crack the json I'm going to be v lazy and just search for the right string
			s := string(body[:])
			log.Println("apps json: " + s)
			var json string = "\"id\":\"/" + app
			index := strings.Index(s, json)

			if index == -1 {
				// Not present, create app. We need to post a json formatted request to the same URL
				// we used above
				//{
				//"id": "xxxxxxxxxx",
				//"cmd": "/tmp/sleep.sh",
				//"cpus": 0.08,
				//"mem": 70,
				//"instances": 1,
				//"container": {
				//    "type": "DOCKER",
				//    "docker": {
				//        "image": "quay.io/rossf7/force12-sleeper:latest",
				//        "network": "BRIDGE"
				//    }
				//}
				//}
				log.Println("create app: " + app)
				var jsonStr string
				jsonStr = "{\"id\":\"xxxxxxxxxx\",\"cmd\":\"/tmp/sleep.sh\",\"cpus\":0.08,\"mem\":70,\"instances\":1,\"container\":{\"type\":\"DOCKER\",\"docker\":{\"image\":\"quay.io/rossf7/force12-sleeper:latest\",\"network\":\"BRIDGE\"}}}"

				jsonStr = strings.Replace(jsonStr, "xxxxxxxxxx", app, 1)

				var a []byte
				copy(a[:], jsonStr)

				resp1, err1 := http.Post(str, "application/json", bytes.NewBuffer(a))
				if err1 != nil {
					// handle error
				}
				defer resp1.Body.Close()

			}
		}

	}
}
