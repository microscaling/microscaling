// Package monitor defines monitors, where we send updates about tasks and performance
package monitor

import (
	"github.com/op/go-logging"

	"github.com/microscaling/microscaling/demand"
)

// Monitor defines the interface for monitors, which receive information about microscaling on a regular basis
type Monitor interface {
	// SendMetrics sends information about the current state of tasks to the monitor
	SendMetrics(tasks *demand.Tasks) (err error)
}

var log = logging.MustGetLogger("mssmonitor")
