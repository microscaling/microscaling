// The demand package defines the interface for demand models
package demand

import (
	"sync"

	"github.com/op/go-logging"
)

type Tasks struct {
	Tasks map[string]Task
	sync.RWMutex
}

type Task struct {
	Demand          int
	Requested       int
	Running         int
	FamilyName      string
	Image           string
	Command         string
	PublishAllPorts bool
	Priority        int
	Env             []string
	MinContainers   int
	MaxContainers   int
	TargetType      string
	Target          int
}

var log = logging.MustGetLogger("mssdemand")

func Exited(tasks *Tasks) (done bool) {
	tasks.RLock()
	defer tasks.RUnlock()

	done = true
	for name, task := range tasks.Tasks {
		if task.Running > 0 {
			done = false
			log.Debugf("Waiting for %s, still %d running, %d requested", name, task.Running, task.Requested)
		}
	}

	return done
}

func ScaleComplete(tasks *Tasks) (done bool) {
	tasks.RLock()
	defer tasks.RUnlock()

	done = true
	for name, task := range tasks.Tasks {
		if task.Running != task.Requested {
			done = false
			log.Debugf("Scale outstanding for %s: %d running, %d requested", name, task.Running, task.Requested)
		}
	}

	return done
}
