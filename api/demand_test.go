package api

import (
	"log"
	"net/http"
	"net/http/httptest"
	// "sync"
	"testing"

	"github.com/force12io/force12/demand"
	"golang.org/x/net/websocket"
)

// var once sync.Once

func TestGetDemand(t *testing.T) {
	var tasks map[string]demand.Task = make(map[string]demand.Task)

	tasks["priority1"] = demand.Task{Demand: 7, Requested: 3, Running: 4}
	tasks["priority2"] = demand.Task{Demand: 3, Requested: 7, Running: 5}

	tests := []struct {
		expJson string
	}{
		{
			expJson: `{
			   "demand": {
			       "tasks": [
			           {
			               "app": "priority1",
			               "demandCount": 7
			           },
			           {
			               "app": "priority2",
			               "demandCount": 3
			           }
			       ]
			   }
			}`,
		},
	}

	for _, test := range tests {

		once.Do(func() {
			http.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
				b := make([]byte, 1000)

				_, err := ws.Read(b)
				if err != nil {
					t.Fatalf("Error reading from web socket %v", err)
				}
				if string(b) != test.expJson {
					log.Printf("Got %v", b)
					t.Fatalf("Unexpected JSON %v", b)
				}
			}))
			server := httptest.NewServer(nil)
			serverAddr = server.Listener.Addr().String()
		})

		baseF12APIUrl = serverAddr
		ws, err := InitWebSocket()
		if err != nil {
			t.Fatal("dialing", err)
		}

		SendMetrics(ws, "hello", tasks)
		ws.Close()
	}
}
