// The demand package defines the interface for demand models
package demand

import (
	"fmt"
	"sort"
	"sync"

	"github.com/op/go-logging"
)

type Tasks struct {
	Tasks []*Task
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


	// Scaling calculation of the ideal number of containers we'd have if there were no other tasks
	IdealContainers int
}

var log = logging.MustGetLogger("mssdemand")

func Exited(tasks *Tasks) (done bool) {
	tasks.RLock()
	defer tasks.RUnlock()

	done = true
	for _, task := range tasks.Tasks {
		if task.Running > 0 {
			done = false
			log.Debugf("Waiting for %s, still %d running, %d requested", task.Name, task.Running, task.Requested)
		}
	}

	return done
}

func ScaleComplete(tasks *Tasks) (done bool) {
	tasks.RLock()
	defer tasks.RUnlock()

	done = true
	for _, task := range tasks.Tasks {
		if task.Running != task.Requested {
			done = false
			log.Debugf("Scale outstanding for %s: %d running, %d requested", task.Name, task.Running, task.Requested)
		}
	}

	return done
}

// implements sort.Interface tasks based on priority
type byPriority []*Task

func (p byPriority) Len() int           { return len(p) }
func (p byPriority) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p byPriority) Less(i, j int) bool { return p[i].Priority < p[j].Priority }

func prioritySort(s []*Task, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(byPriority(s)))
	} else {
		sort.Sort(byPriority(s))
	}
}

func (t *Tasks) GetTask(name string) (task *Task, err error) {
	for _, task := range t.Tasks {
		if task.Name == name {
			return task, nil
		}
	}

	return nil, fmt.Errorf("No Task with name %s", name)
}
