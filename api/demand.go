// API between Force12 agent and server
package api

import (
	"log"

	"golang.org/x/net/websocket"
)

type DemandPayload struct {
	Demand DemandUpdate `json:"demand"`
}

type DemandUpdate struct {
	Tasks []TaskDemand `json:"tasks"`
}

type TaskDemand struct {
	App         string `json:"app"`
	DemandCount int    `json:"demandCount"`
}

func Listen(ws *websocket.Conn, demandUpdate chan []TaskDemand) error {
	for {
		select {
		default:
			var dp DemandPayload
			err := websocket.JSON.Receive(ws, &dp)
			if err != nil {
				log.Printf("Error reading from web socket: %v", err)
				return err
			} else {
				log.Printf("Received demand %v", dp)
				demandUpdate <- dp.Demand.Tasks
			}
		}
	}
}
