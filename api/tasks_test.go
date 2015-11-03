package api

import (
	"encoding/json"
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

	_ = tasksFromResponse(b, tasks)

	p1 := tasks["priority1"]
	if p1.Demand != 99 {
		t.Fatalf("Bad demand for p1: %v", p1)
	}
	if p1.Requested != 5 {
		t.Fatalf("Requested got modified for p1: %v", p1)
	}
	p2 := tasks["another_app"]
	if p2.Demand != 100 {
		t.Fatalf("Bad demand for p2")
	}
	// Demand should be unchanged for this one
	p3 := tasks["p3"]
	if p3.Demand != 9 {
		t.Fatalf("Bad demand for p3")
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
		tasks   map[string]demand.Task
	}{
		{
			expUrl:  "/tasks/hello",
			json:    "",
			success: false,
			tasks: map[string]demand.Task{
				"priority1": demand.Task{
					Image:     "firstimage",
					Demand:    8,
					Requested: 3,
					Running:   4,
				},
				"priority2": demand.Task{
					Image:     "anotherimage",
					Command:   "do this",
					Demand:    0,
					Requested: 7,
					Running:   5,
				},
			},
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
			tasks: map[string]demand.Task{
				"priority1": demand.Task{
					Image:     "firstimage",
					Demand:    7,
					Requested: 3,
					Running:   4,
				},
				"priority2": demand.Task{
					Image:     "anotherimage",
					Command:   "do this",
					Demand:    3,
					Requested: 7,
					Running:   5,
				},
			},
		},
	}

	for _, test := range tests {
		server := apitest.DoTestGetJson(t, test.expUrl, test.success, test.json)
		defer server.Close()

		baseF12APIUrl = server.URL
		_ = GetTasks("hello", tasks)
		baseF12APIUrl = getBaseF12APIUrl()

		checkReturnedTasks(t, test.tasks, tasks)
	}
}
