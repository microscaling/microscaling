// The demand package defines the interface for demand models
package demand

import (
	"sync"
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
}
