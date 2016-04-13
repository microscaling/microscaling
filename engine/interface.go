package engine

import (
	"github.com/op/go-logging"

	"github.com/microscaling/microscaling/demand"
)

type Engine interface {
	// Engine is responsible for setting up the new Demand in tasks. If demand has changed, send on demandUpdate
	GetDemand(tasks *demand.Tasks, demandUpdate chan struct{})

	// When the engine has cleaned itself up, it must close this demandUpdate channel
	StopDemand(demandUpdate chan struct{})
}

var log = logging.MustGetLogger("mssengine")
