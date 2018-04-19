package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/microscaling/microscaling/demand"
)

func TestGetAppsDecode(t *testing.T) {
	var bb = []byte(`{"image": "my image", "command": "do it"}`)
	var d = DockerAppConfig{}
	_ = json.Unmarshal(bb, &d)
	if d.Image != "my image" {
		t.Fatalf("Didn't decode image")
	}

	if d.Command != "do it" {
		t.Fatalf("Didn't decode command")
	}

	// var response string = `"apps": [{"name":"priority1","appType":"Docker","config":{"image":"force12io/priority-1:latest","command":"/run.sh"}},{"name":"priority2","type":"Docker","config":{"image":"force12io/priority-2:latest","command":"/run.sh"}}]`
	var response = `{"apps" : [{"name":"priority1", "config":{"image":"microscaling/priority-1:latest","command":"/run.sh"}},{"name":"priority2","appType":"Docker","config":{"image":"microscaling/priority-2:latest","command":"/run.sh"}}]}`
	var b = []byte(response)

	var a AppsMessage
	err := json.Unmarshal(b, &a)
	if err != nil {
		t.Fatalf("Error decoding apps message: %v", err)
	}

	var apps []*demand.Task
	apps, _, _ = appsFromResponse(b)

	var p1, p2 *demand.Task
	for _, task := range apps {
		switch task.Name {
		case "priority1":
			p1 = task
		case "priority2":
			p2 = task
		}
	}
	if p1.Image != "microscaling/priority-1:latest" {
		t.Fatalf("Bad image %s", p1.Image)
	}
	if p2.Image != "microscaling/priority-2:latest" {
		t.Fatalf("Bad image %s", p2.Image)
	}
}

func TestGetApps(t *testing.T) {
	// Needed to create SQS metric
	os.Setenv("AWS_REGION", "us-east-1")

	tests := []struct {
		expURL      string
		json        string
		success     bool
		tasks       map[string]demand.Task
		metricTypes map[string]string
		targetTypes map[string]string
	}{
		{
			expURL: "/apps/hello",
			json: `{"apps": [
			      {
			          "name": "priority1",
			          "appType": "Docker",
			          "ruleType": "Queue",
			          "metricType": "NSQ",
			          "config": {
			              "image": "firstimage",
			              "topicName": "test",
			              "channelName": "test"
			          }
			      },
			      {
			          "name": "priority2",
			          "appType": "Docker",
			          "ruleType": "SimpleQueue",
			          "metricType": "SQS",
			          "config": {
			              "image": "anotherimage",
			              "command": "do this",
			              "queueURL": "https://sqs.us-east-1.amazonaws.com/12345/test"
			          }
			      }
			]}`,
			success: true,
			tasks: map[string]demand.Task{
				"priority1": demand.Task{
					Image: "firstimage",
				},
				"priority2": demand.Task{
					Image:   "anotherimage",
					Command: "do this",
				},
			},
			metricTypes: map[string]string{
				"priority1": "*metric.NSQMetric",
				"priority2": "*metric.SQSMetric",
			},
			targetTypes: map[string]string{
				"priority1": "*target.QueueLengthTarget",
				"priority2": "*target.SimpleQueueLengthTarget",
			},
		},
		{
			expURL:  "/apps/hello",
			json:    "",
			success: false,
			tasks:   map[string]demand.Task{},
		},
	}

	for _, test := range tests {
		server := DoTestGetJSON(t, test.expURL, test.success, test.json)
		defer server.Close()

		baseAPIUrl := strings.Replace(server.URL, "http://", "", 1)
		returnedTasks, _, err := GetApps(baseAPIUrl, "hello")

		if test.success {
			CheckReturnedTasks(t, test.tasks, returnedTasks)
			CheckMetricTypes(t, test.metricTypes, returnedTasks)
			CheckTargetTypes(t, test.targetTypes, returnedTasks)

			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

		} else {
			if err == nil {
				t.Fatalf("Expected an error")
			}
		}
	}
}

// Utility for checking GET requests
func DoTestGetJSON(t *testing.T, expURL string, success bool, testJSON string) (server *httptest.Server) {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expURL {
			t.Fatalf("Expected %s, have %s", expURL, r.URL.Path)
		}

		if r.Method != "GET" {
			t.Fatalf("expected GET, have %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Content type not as expected, have %s", r.Header.Get("Content-Type"))
		}

		if success {
			w.Write([]byte(testJSON))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))

	return server
}

// Utility for checking that tasks are updated to be what we expect
func CheckReturnedTasks(t *testing.T, tasks map[string]demand.Task, returnedTasks []*demand.Task) {
	for _, rt := range returnedTasks {
		tt, ok := tasks[rt.Name]
		if !ok {
			t.Fatalf("Unexpected app name %v", rt.Name)
		}

		if tt.Image != rt.Image {
			t.Fatalf("Image: expected %s got %s", tt.Image, rt.Image)
		}
		if tt.Command != rt.Command {
			t.Fatalf("Command: expected %s got %s", tt.Command, rt.Command)
		}
		if tt.Demand != rt.Demand {
			t.Fatalf("Demand: expected %d got %d", tt.Demand, rt.Demand)
		}
		if tt.Requested != rt.Requested {
			t.Fatalf("Requested: expected %d got %d", tt.Requested, rt.Requested)
		}
		if tt.Running != rt.Running {
			t.Fatalf("Requested: expected %d got %d", tt.Requested, rt.Requested)
		}
	}
}

// CheckMetricTypes checks the correct scaling metric was selected
func CheckMetricTypes(t *testing.T, metricTypes map[string]string, returnedTasks []*demand.Task) {
	for _, rt := range returnedTasks {
		typeName := reflect.TypeOf(rt.Metric).String()

		if typeName != metricTypes[rt.Name] {
			t.Fatalf("MetricType: expected %s got %s", metricTypes[rt.Name], typeName)
		}
	}
}

// CheckTargetTypes checks the correct scaling target was selected
func CheckTargetTypes(t *testing.T, targetTypes map[string]string, returnedTasks []*demand.Task) {
	for _, rt := range returnedTasks {
		typeName := reflect.TypeOf(rt.Target).String()

		if typeName != targetTypes[rt.Name] {
			t.Fatalf("TargetType: expected %s got %s", targetTypes[rt.Name], typeName)
		}
	}
}
