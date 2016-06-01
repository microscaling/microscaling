package monitor

import (
	"net/http/httptest"
	"testing"

	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/utils"
	"golang.org/x/net/websocket"
)

func testServerMetrics(ws *websocket.Conn) {
	// Contents are tested elsewhere, this just checks we can read something off the socket
	var b []byte
	b = make([]byte, 1000)
	ws.Read(b)
}

func TestServerMonitor(t *testing.T) {
	var tasks demand.Tasks
	tasks.Tasks = make([]*demand.Task, 2)

	tasks.Tasks[0] = &demand.Task{Name: "priority1", Demand: 8, Requested: 3, Running: 4}
	tasks.Tasks[1] = &demand.Task{Name: "priority2", Demand: 2, Requested: 7, Running: 5}

	server := httptest.NewServer(websocket.Handler(testServerMetrics))
	serverAddr := server.Listener.Addr().String()

	ws, err := utils.InitWebSocket(serverAddr)
	if err != nil {
		t.Fatal("dialing", err)
	}

	s := NewServerMonitor(ws, "hello")
	if s.userID != "hello" {
		t.Fatal("Didn't set userID")
	}
	s.SendMetrics(&tasks)

	ws.Close()
	server.Close()
}
