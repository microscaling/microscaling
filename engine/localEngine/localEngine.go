package localEngine

import (
	"sync"
	"time"

	"github.com/op/go-logging"

	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/engine"
)

const constGetDemandSleep = 500

// LocalEngine calculates demand locally
type LocalEngine struct {
}

// compile-time assert that we implement the right interface
var _ engine.Engine = (*LocalEngine)(nil)

var log = logging.MustGetLogger("mssengine")

// NewEngine initializes the local engine
func NewEngine() *LocalEngine {
	de := LocalEngine{}
	return &de
}

func (de *LocalEngine) GetDemand(tasks *demand.Tasks, demandUpdate chan struct{}) {
	var gettingMetrics sync.WaitGroup

	// In this we need to collect the metrics, calculate demand, and trigger a demand update
	demandTimeout := time.NewTicker(constGetDemandSleep * time.Millisecond)
	for _ = range demandTimeout.C {
		tasks.Lock()
		log.Debug("Getting demand")

		for _, task := range tasks.Tasks {
			go func() {
				gettingMetrics.Add(1)
				defer gettingMetrics.Done()
				task.Metric.UpdateCurrent()
			}
		}

		gettingMetrics.Wait()

		demandChanged := ScalingCalculation(tasks)

		tasks.Unlock()
		if demandChanged {
			demandUpdate <- struct{}{}
		}
	}
}

// StopDemand is called when we want to shut down
func (de *LocalEngine) StopDemand(demandUpdate chan struct{}) {
	close(demandUpdate)
}
