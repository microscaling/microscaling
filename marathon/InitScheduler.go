package marathon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type app struct {
	Id        string    `json:"id"`
	Cmd       string    `json:"cmd"`
	Cpus      float64   `json:"cpus"`
	Mem       int       `json:"mem"`
	Instances int       `json:"instances"`
	Container container `json:"container"`
}

type apps struct {
	Apps []app `json:"apps"`
}

type docker struct {
	Image   string `json:"image"`
	Network string `json:"network"`
}

type container struct {
	Type   string `json:"type"`
	Docker docker `json:"docker"`
}

// InitScheduler ensures our app is started on the cluster.
//
// appId - string identifier of a container
func (m *MarathonScheduler) InitScheduler(appId string) error {
	// Ask Marathon which apps are running
	resp, err := http.Get(m.baseMarathonUrl)
	defer resp.Body.Close()
	if err != nil {
		// handle error
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Error response from marathon, %d %s", resp.StatusCode, resp.Status)
	}

	payload := apps{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		return fmt.Errorf("Failed to decode json response. %v", err)
	}

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
	haveApp := false
	for _, app := range payload.Apps {
		if app.Id == appId {
			haveApp = true
			break
		}
	}

	if !haveApp {
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
		log.Printf("create app: %s", appId)

		app := app{
			Id:        appId,
			Cmd:       "/tmp/sleep.sh",
			Cpus:      0.08,
			Mem:       70,
			Instances: 1,
			Container: container{
				Type: "DOCKER",
				Docker: docker{
					Image:   "quay.io/rossf7/force12-sleeper:latest",
					Network: "BRIDGE",
				},
			},
		}

		w := &bytes.Buffer{}
		encoder := json.NewEncoder(w)
		err := encoder.Encode(&app)
		if err != nil {
			return fmt.Errorf("Failed to encode json. %v", err)
		}

		resp, err := http.Post(m.baseMarathonUrl, "application/json", w)
		if err != nil {
			// handle error
			return fmt.Errorf("Failed to create app, err %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode > 299 {
			return fmt.Errorf("Error response from marathon, %d %s", resp.StatusCode, resp.Status)
		}

	}
	return nil
}
