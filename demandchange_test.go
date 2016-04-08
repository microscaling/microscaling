package main

import (
	"log"
	"reflect"
	"testing"

	"github.com/microscaling/microscaling/api"
	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/scheduler/toy"
)

func TestHandleDemandChange(t *testing.T) {
	tasks = make(map[string]demand.Task)
	tasks["priority1"] = demand.Task{
		FamilyName: "p1family",
		Demand:     4,
		Requested:  0,
	}

	tasks["priority2"] = demand.Task{
		FamilyName: "p2family",
		Demand:     3,
		Requested:  0,
	}

	s := toy_scheduler.NewScheduler()

	tests := []struct {
		td       []api.TaskDemand
		newtasks map[string]demand.Task
	}{
		{
			td: []api.TaskDemand{
				{
					App:         "priority1",
					DemandCount: 5,
				},
			},
			newtasks: map[string]demand.Task{
				"priority1": {FamilyName: "p1family", Demand: 5, Requested: 5},
				"priority2": {FamilyName: "p2family", Demand: 3, Requested: 3},
			},
		},
		{
			// We just ignore any tasks that we didn't know about
			td: []api.TaskDemand{
				{
					App:         "priority1",
					DemandCount: 5,
				}, {
					App:         "priority3",
					DemandCount: 5,
				},
			},
			newtasks: map[string]demand.Task{
				"priority1": {FamilyName: "p1family", Demand: 5, Requested: 5},
				"priority2": {FamilyName: "p2family", Demand: 3, Requested: 3},
			},
		},
	}

	for _, test := range tests {
		err := handleDemandChange(test.td, s, tasks)
		if err != nil {
			t.Fatalf("handleDemandChange failed")
		}
		log.Println(tasks)
		if !reflect.DeepEqual(tasks, test.newtasks) {
			t.Fatalf("Expected %v tasks, got %v", test.newtasks, tasks)
		}
	}
}
