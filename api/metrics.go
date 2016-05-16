// API between Microscaling agent and server
package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/microscaling/microscaling/demand"
	"golang.org/x/net/websocket"
)

// sendMetrics sends the current state of tasks to the API
func SendMetrics(ws *websocket.Conn, userID string, tasks *demand.Tasks) error {
	var err error = nil
	var index int = 0

	metrics := metrics{
		Tasks:     make([]taskMetrics, len(tasks.Tasks)),
		CreatedAt: time.Now().Unix(),
	}

	tasks.Lock()
	for _, task := range tasks.Tasks {
		metrics.Tasks[index] = taskMetrics{App: task.Name, RunningCount: task.Running, PendingCount: task.Requested}

		if task.Metric != nil {
			metrics.Tasks[index].Metric = task.Metric.Current()
		}

		index++
	}
	tasks.Unlock()

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
