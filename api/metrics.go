// API between Force12 agent and server
package api

import (
	// "bytes"
	"encoding/json"
	"fmt"
	// "log"
	// "net/http"
	"time"

	"github.com/force12io/force12/demand"
	"golang.org/x/net/websocket"
)

type metricsPayload struct {
	User      string       `json:"user"`
	CreatedAt int64        `json:"createdAt"`
	Tasks     []appMetrics `json:"tasks"`
}

type appMetrics struct {
	App          string `json:"app"`
	RunningCount int    `json:"runningCount"`
	PendingCount int    `json: "pendingCount"`
}

// sendMetrics sends the current state of tasks to the F12 API
func SendMetrics(ws *websocket.Conn, userID string, tasks map[string]demand.Task) error {
	var err error = nil
	var index int = 0

	// url := baseF12APIUrl + "/metrics/" + userID

	payload := metricsPayload{
		User:      userID,
		CreatedAt: time.Now().Unix(),
		Tasks:     make([]appMetrics, len(tasks)),
	}

	for name, task := range tasks {
		payload.Tasks[index] = appMetrics{App: name, RunningCount: task.Running, PendingCount: task.Requested}
		index++
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to encode API json. %v", err)
	}

	_, err = ws.Write(b)
	if err != nil {
		return fmt.Errorf("Failed to send metrics: %v", err)
	}

	return err
}
