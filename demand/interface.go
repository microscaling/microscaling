// The demand package defines the interface for demand models
package demand

import (
	"sync"

	"github.com/op/go-logging"

	"github.com/microscaling/microscaling/metric"
	"github.com/microscaling/microscaling/target"
)

type Tasks struct {
	Tasks         []*Task
	MaxContainers int
	sync.RWMutex
}

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
