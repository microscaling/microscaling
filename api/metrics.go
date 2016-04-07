// API between Microscaling agent and server
package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/microscaling/microscaling/demand"
	"golang.org/x/net/websocket"
)

type metricsPayload struct {
	User    string  `json:"user"`
	Metrics metrics `json:"metrics"`
}

type metrics struct {
	Tasks     []taskMetrics `json:"tasks"`
	CreatedAt int64         `json:"createdAt"`
}

type taskMetrics struct {
	App          string `json:"app"`
	RunningCount int    `json:"runningCount"`
	PendingCount int    `json:"pendingCount"`
}

// sendMetrics sends the current state of tasks to the API
func SendMetrics(ws *websocket.Conn, userID string, tasks map[string]demand.Task) error {
	var err error = nil
	var index int = 0

	metrics := metrics{
		Tasks:     make([]taskMetrics, len(tasks)),
		CreatedAt: time.Now().Unix(),
	}

	for name, task := range tasks {
		metrics.Tasks[index] = taskMetrics{App: name, RunningCount: task.Running, PendingCount: task.Requested}
		index++
	}

	payload := metricsPayload{
		User:    userID,
		Metrics: metrics,
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
