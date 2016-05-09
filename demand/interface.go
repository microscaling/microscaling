// Package demand defines Tasks
package demand

import (
	"sync"

	"github.com/op/go-logging"

	"github.com/microscaling/microscaling/metric"
	"github.com/microscaling/microscaling/target"
)

// Tasks is a list of tasks, with a global lock. Global config about tasks can go here too.
type Tasks struct {
	Tasks         []*Task
	MaxContainers int
	sync.RWMutex
}

// Task describes an app (or you might want to call it a service, or a container). It has all the info
// for starting / stopping an instance of a task, scaling config & params, the target and metric we use
// for this task, and state information about the number of tasks.
type Task struct {
	// Name
	Name string

	// Scheduler
	Demand    int
	Requested int
	Running   int

	// Container config info
	FamilyName      string
	Image           string
	Command         string
	PublishAllPorts bool
	NetworkMode     string
	Env             []string

	// Scaling config
	IsScalable    bool
	Priority      int
	MaxDelta      int
	MinContainers int
	MaxContainers int

	// The target we're aiming for
	Target target.Target

	// Measurements
	Metric metric.Metric

	// Scaling calculation of the ideal number of containers we'd have if there were no other tasks
	IdealContainers int
}

var log = logging.MustGetLogger("mssdemand")
