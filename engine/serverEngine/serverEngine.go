package serverEngine

import (
	"golang.org/x/net/websocket"

	"github.com/op/go-logging"

	"github.com/microscaling/microscaling/api"
	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/engine"
)

// ServerEngine retrieves demand from the server where it has been calculated
type ServerEngine struct {
	ws        *websocket.Conn
	closedown chan struct{}
}

// compile-time assert that we implement the right interface
var _ engine.Engine = (*ServerEngine)(nil)

var log = logging.MustGetLogger("mssengine")

// NewEngine initializes a new ServerEngine
func NewEngine(ws *websocket.Conn) *ServerEngine {
	de := ServerEngine{
		ws:        ws,
		closedown: make(chan struct{}, 1),
	}
	return &de
}

func updateTasks(dp api.DemandPayload, tasks *demand.Tasks) (demandChanged bool) {
	demandChanged = false
	tasks.Lock()
	defer tasks.Unlock()

	for _, taskFromServer := range dp.Demand.Tasks {
		name := taskFromServer.App

		if existingTask, err := tasks.GetTask(name); err == nil {
			if existingTask.Demand != taskFromServer.DemandCount {
				demandChanged = true
			}
			existingTask.Demand = taskFromServer.DemandCount
		}
	}
	return demandChanged
}

// GetDemand collects the metrics, calculates demand, and sends demandUpdate messages on the API
func (de *ServerEngine) GetDemand(tasks *demand.Tasks, demandUpdate chan struct{}) {
	var demandChanged bool
	var dp api.DemandPayload
	var closedown = false
	for {
		select {
		// We can't just close the websocket on closedown as we (may) still want to send metrics on it
		case <-de.closedown:
			closedown = true
		default:
			err := websocket.JSON.Receive(de.ws, &dp)
			if err != nil {
				log.Errorf("Error reading from web socket: %v", err)
				break
			}

			if closedown {
				close(demandUpdate)
				log.Debug("[server] channel closed")
				break
			}

			log.Debugf("Received demand %v", dp)
			demandChanged = updateTasks(dp, tasks)
			if demandChanged {
				demandUpdate <- struct{}{}
			}
		}
	}
}

// StopDemand should be called when we want to shut down. We signal the ServerEngine to close the websocket.
func (de *ServerEngine) StopDemand(demandUpdate chan struct{}) {
	log.Debug("[server] stop receiving demand")
	de.closedown <- struct{}{}
}
