package monitor

import (
	"golang.org/x/net/websocket"

	"github.com/microscaling/microscaling/api"
	"github.com/microscaling/microscaling/demand"
)

// ServerMonitor receives metrics over the Microscaling API on the websocket
type ServerMonitor struct {
	ws     *websocket.Conn
	userID string
}

// NewServerMonitor returns a new monitor that uses the Microscaling API to send metrics about running tasks and demand
func NewServerMonitor(ws *websocket.Conn, userID string) *ServerMonitor {
	return &ServerMonitor{
		ws:     ws,
		userID: userID,
	}
}

// SendMetrics uses the Microscaling API to send metrics about all the running tasks and demand
func (m *ServerMonitor) SendMetrics(tasks *demand.Tasks) (err error) {
	err = api.SendMetrics(m.ws, m.userID, tasks)
	return
}
