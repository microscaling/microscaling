package compose

import (
	"testing"

	"bitbucket.org/force12io/force12-scheduler/demand"
)

func TestComposeScheduler(t *testing.T) {
	c := NewScheduler()
	var task demand.Task
	task.Demand = 5

	c.InitScheduler("anything", &task)
}
