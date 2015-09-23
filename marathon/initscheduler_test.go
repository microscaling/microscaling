package marathon

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const initSchedRsp = `{
	"apps": [
	 {
	     "id": "myapp",
	     "cmd": "/tmp/sleep.sh",
	     "args": null,
	     "user": null,
	     "env": {},
	     "instances": 2,
	     "cpus": 0.08,
	     "mem": 70,
	     "disk": 0,
	     "executor": "",
	     "constraints": [],
	     "uris": [],
	     "storeUrls": [],
	     "ports": [
	         10000
	     ],
	     "requirePorts": false,
	     "backoffSeconds": 1,
	     "backoffFactor": 1.15,
	     "maxLaunchDelaySeconds": 3600,
	     "container": {
	         "type": "DOCKER",
	         "volumes": [],
	         "docker": {
	             "image": "quay.io/rossf7/force12-sleeper:latest",
	             "network": "BRIDGE",
	             "privileged": false,
	             "parameters": [],
	             "forcePullImage": false
	         }
	     },
	     "healthChecks": [],
	     "dependencies": [],
	     "upgradeStrategy": {
	         "minimumHealthCapacity": 1,
	         "maximumOverCapacity": 1
	     },
	     "labels": {},
	     "acceptedResourceRoles": null,
	     "version": "2015-07-29T18:03:35.688Z",
	     "tasksStaged": 0,
	     "tasksRunning": 2,
	     "tasksHealthy": 0,
	     "tasksUnhealthy": 0,
	     "deployments": []
	 }
	]
}`

func TestInitScheduler(t *testing.T) {

	tests := []struct {
		appId          string
		expPost        bool
		expGet         bool
		postStatus     int
		getStatus      int
		getPayload     string
		expPostPayload string
		expError       bool
	}{
		{
			appId:      "myapp",
			expGet:     true,
			getStatus:  200,
			getPayload: initSchedRsp,
		},
		{
			appId:      "myotherapp",
			expGet:     true,
			getStatus:  200,
			getPayload: initSchedRsp,
			expPost:    true,
			postStatus: 200,
			expPostPayload: `{"id":"myotherapp","cmd":"/tmp/sleep.sh","cpus":0.08,"mem":70,"instances":1,"container":{"type":"DOCKER","docker":{"image":"quay.io/rossf7/force12-sleeper:latest","network":"BRIDGE"}}}
`,
		},
		{
			appId:      "myapp",
			expGet:     true,
			getStatus:  500,
			getPayload: `{"error": "oh no!"}`,
			expError:   true,
		},
		{
			appId:      "myotherapp",
			expGet:     true,
			getStatus:  200,
			getPayload: initSchedRsp,
			expPost:    true,
			postStatus: 500,
			expPostPayload: `{"id":"myotherapp","cmd":"/tmp/sleep.sh","cpus":0.08,"mem":70,"instances":1,"container":{"type":"DOCKER","docker":{"image":"quay.io/rossf7/force12-sleeper:latest","network":"BRIDGE"}}}
`,
			expError: true,
		},
	}

	for _, test := range tests {
		var havePost, haveGet bool
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				t.Fatalf("Expected root path, have %s", r.URL.Path)
			}

			switch r.Method {
			case "POST":
				havePost = true
				if r.Header.Get("Content-Type") != "application/json" {
					t.Fatalf("Content type not as expected, have %s", r.Header.Get("Content-Type"))
				}
				data, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("Failed to read post body. %v", err)
				}
				if string(data) != test.expPostPayload {
					t.Fatalf("post payload not as expected have %s, expected **%s**", string(data), test.expPostPayload)
				}
				w.WriteHeader(test.postStatus)
			case "GET":
				haveGet = true
				h := w.Header()
				h.Set("Content-Type", "application/json")
				w.WriteHeader(test.getStatus)
				fmt.Fprintln(w, test.getPayload)
			default:
				t.Fatalf("Unexpected method %s", r.Method)
			}

		}))
		defer server.Close()

		m := NewScheduler()
		m.baseMarathonUrl = server.URL

		err := m.InitScheduler(test.appId)
		if err != nil {
			if !test.expError {
				t.Fatalf("InitScheduler returned an error %v", err)
			}
		} else if test.expError {
			t.Fatalf("Expected an error, did not get one")
		}

		if havePost != test.expPost {
			t.Fatalf("Post expectation - have post %t", havePost)
		}
		if haveGet != test.expGet {
			t.Fatalf("Get expectation - have get %t", haveGet)
		}
	}
}
