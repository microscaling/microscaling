package serverEngine

import (
	"golang.org/x/net/websocket"

	"github.com/op/go-logging"

	"github.com/microscaling/microscaling/api"
	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/engine"
)

type ServerEngine struct {
	ws        *websocket.Conn
	closedown chan struct{}
}

// compile-time assert that we implement the right interface
var _ engine.Engine = (*ServerEngine)(nil)

var log = logging.MustGetLogger("mssengine")

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

func (de *ServerEngine) GetDemand(tasks *demand.Tasks, demandUpdate chan struct{}) {
	// In this engine the Server collects the metrics, calculates demand, and sends demandUpdate messages on the API	api.Listen(de.ws, tasks, demandUpdate)
	var demandChanged bool
	var dp api.DemandPayload
	var closedown bool = false
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

func (de *ServerEngine) StopDemand(demandUpdate chan struct{}) {
	log.Debug("[server] stop receiving demand")
	de.closedown <- struct{}{}
}
