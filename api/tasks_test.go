package api

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/force12io/force12/api/apitest"
	"github.com/force12io/force12/demand"
)

func TestGetTasksDecode(t *testing.T) {
	var response string = `[{"app":"priority1", "demandCount": 99},{"app":"another_app","demandCount": 100}]`
	var b = []byte(response)

	var tt []TaskDescription
	_ = json.Unmarshal(b, &tt)

	var tasks map[string]demand.Task = make(map[string]demand.Task, 3)
	t1 := demand.Task{
		Demand:    3,
		Requested: 5,
	}
	t2 := demand.Task{}
	tasks["priority1"] = t1
	tasks["another_app"] = t2
	tasks["p3"] = demand.Task{
		Demand: 9,
	}
}

func TestGetTasks(t *testing.T) {
	var tasks map[string]demand.Task = make(map[string]demand.Task)

	tasks["priority1"] = demand.Task{Image: "firstimage", Demand: 8, Requested: 3, Running: 4}
	tasks["priority2"] = demand.Task{Image: "anotherimage", Command: "do this", Demand: 0, Requested: 7, Running: 5}

	tests := []struct {
		expUrl  string
		json    string
		success bool
		td      []TaskDescription
	}{
		{
			expUrl:  "/tasks/hello",
			json:    "",
			success: false,
			td:      []TaskDescription{},
		}, {
			expUrl: "/tasks/hello",
			json: `[
			    {
			        "app": "priority1",
			        "demandCount": 7
			    },
			    {
			        "app": "priority2",
			        "demandCount": 3
			    }
			]`,
			success: true,
			td: []TaskDescription{
				TaskDescription{
					App:         "priority1",
					DemandCount: 7,
				},
				TaskDescription{
					App:         "priority2",
					DemandCount: 3,
				},
			},
		},
	}

	for number, test := range tests {
		server := apitest.DoTestGetJson(t, test.expUrl, test.success, test.json)
		defer server.Close()

		baseF12APIUrl = server.URL
		td, err := GetTasks("hello")
		baseF12APIUrl = getBaseF12APIUrl()

		if test.success {
			if err != nil {
				t.Fatalf("Didn't expect failure: v%", err)
			}
			if !reflect.DeepEqual(td, test.td) {
				t.Fatalf("Task descriptions not equal: %v | %v", td, test.td)
			}
		} else {
			if err == nil {
				t.Fatalf("Test %d was supposed to fail", number)
			}
		}
	}
}
