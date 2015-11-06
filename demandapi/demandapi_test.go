package demandapi

import (
	"testing"

	"github.com/force12io/force12/api"
	"github.com/force12io/force12/api/apitest"
	"github.com/force12io/force12/demand"
)

func TestDemandApi(t *testing.T) {
	var demandChanged bool
	var tasks map[string]demand.Task = make(map[string]demand.Task)

	tasks["priority1"] = demand.Task{Image: "firstimage", Demand: 8, Requested: 3, Running: 4}
	tasks["priority2"] = demand.Task{Image: "anotherimage", Command: "do this", Demand: 0, Requested: 7, Running: 5}

	d := NewDemandModel("hello")

	tests := []struct {
		expUrl        string
		json          string
		success       bool
		demandChanged bool
		tasks         map[string]demand.Task
	}{
		{
			expUrl:        "/tasks/hello",
			json:          "",
			success:       false,
			demandChanged: false,
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
			        "demandCount": 8
			    },
			    {
			        "app": "priority2",
			        "demandCount": 0
			    }
			]`,
			success:       true,
			demandChanged: false,
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
			success:       true,
			demandChanged: true,
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

		api.SetBaseF12APIUrl(server.URL)
		demandChanged, _ = d.Update(tasks)
		api.SetBaseF12APIUrl(api.GetBaseF12APIUrl())

		if test.demandChanged != demandChanged {
			t.Fatalf("Demand changed: %v expected %v", demandChanged, test.demandChanged)
		}
		apitest.CheckReturnedTasks(t, test.tasks, tasks)
	}
}
