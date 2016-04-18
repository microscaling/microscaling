package api

import (
	"encoding/json"
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
	var response string = `{"apps" : [{"name":"priority1", "config":{"image":"microscaling/priority-1:latest","command":"/run.sh"}},{"name":"priority2","appType":"Docker","config":{"image":"microscaling/priority-2:latest","command":"/run.sh"}}]}`
	var b = []byte(response)

	var a AppsMessage
	err := json.Unmarshal(b, &a)
	if err != nil {
		t.Fatalf("Error decoding apps message: %v", err)
	}

	var apps map[string]demand.Task
	apps, _, _ = appsFromResponse(b)

	p1 := apps["priority1"]
	if p1.Image != "microscaling/priority-1:latest" {
		t.Fatalf("Bad image %s", p1.Image)
	}
	p2 := apps["priority2"]
	if p2.Image != "microscaling/priority-2:latest" {
		t.Fatalf("Bad image %s", p2.Image)
	}
}

func TestGetApps(t *testing.T) {
	tests := []struct {
		expUrl  string
		json    string
		success bool
		tasks   map[string]demand.Task
	}{
		{
			expUrl: "/apps/hello",
			json: `{"apps": [
			      {
			          "name": "priority1",
			          "appType": "Docker",
			          "config": {
			              "image": "firstimage"
			          }
			      },
			      {
			          "name": "priority2",
			          "appType": "Docker",
			          "config": {
			              "image": "anotherimage",
			              "command": "do this"
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
		},
		{
			expUrl:  "/apps/hello",
			json:    "",
			success: false,
			tasks:   map[string]demand.Task{},
		},
	}

	for _, test := range tests {
		server := DoTestGetJson(t, test.expUrl, test.success, test.json)
		defer server.Close()

		baseAPIUrl = strings.Replace(server.URL, "http://", "", 1)
		returned_tasks, _, err := GetApps("hello")
		baseAPIUrl = GetBaseAPIUrl()

		if test.success {
			CheckReturnedTasks(t, test.tasks, returned_tasks)

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
