package demand

import (
	"fmt"
	"sort"
)

// Exited returns whether tasks have all drained down to 0 so we can quit microscaling
func (tasks *Tasks) Exited() (done bool) {
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

// CheckCapacity returns number of containers we have space for
func (tasks *Tasks) CheckCapacity() int {
	// TODO!! For now we are simply going to say there is a maximum total number of containers this deployment can handle
	// TODO!! It should really look at the available CPU / mem / bw in / out
	totalRequested := 0
	for _, t := range tasks.Tasks {
		totalRequested += t.Requested
	}

	return tasks.MaxContainers - totalRequested
}

// implements sort.Interface tasks based on priority
type byPriority []*Task

func (p byPriority) Len() int           { return len(p) }
func (p byPriority) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p byPriority) Less(i, j int) bool { return p[i].Priority < p[j].Priority }

// PrioritySort reorders tasks in priority order (or reverse priority order)
func (tasks *Tasks) PrioritySort(reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(byPriority(tasks.Tasks)))
	} else {
		sort.Sort(byPriority(tasks.Tasks))
	}
}

// GetTask returns the task identified by name
func (tasks *Tasks) GetTask(name string) (task *Task, err error) {
	for _, task := range tasks.Tasks {
		if task.Name == name {
			return task, nil
		}
	}

	return nil, fmt.Errorf("No Task with name %s", name)
}
