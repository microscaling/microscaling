package api

import (
	// "encoding/json"
	// "io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/force12io/force12/demand"
	"golang.org/x/net/websocket"
)

var serverAddr string
var once sync.Once

func testServer(ws *websocket.Conn) {
	log.Printf("Received something")
}

func startServer() {
	http.Handle("/", websocket.Handler(testServer))
	server := httptest.NewServer(nil)
	serverAddr = server.Listener.Addr().String()
	log.Print("Test WebSocket server listening on ", serverAddr)
}

func TestInitWebSocket(t *testing.T) {
	once.Do(startServer)

	baseF12APIUrl = serverAddr
	ws, err := InitWebSocket()
	if err != nil {
		t.Fatal("dialing", err)
	}

	msg := []byte("hello, world\n")
	if _, err := ws.Write(msg); err != nil {
		t.Errorf("Write: %v", err)
	}
	ws.Close()
}

func TestSendMetrics(t *testing.T) {
	var tasks map[string]demand.Task = make(map[string]demand.Task)

	tasks["priority1"] = demand.Task{Demand: 8, Requested: 3, Running: 4}
	tasks["priority2"] = demand.Task{Demand: 2, Requested: 7, Running: 5}

	tests := []struct {
		expJson string
	}{
		{
			expJson: `{
			   "user": "5k4ek",
			   "createdAt": 1435071103,
			   "metrics": {
			       "tasks": [
			           {
			               "app": "priority1",
			               "runningCount": 4,
			               "pendingCount": 3
			           },
			           {
			               "app": "priority2",
			               "runningCount": 5,
			               "pendingCount": 7
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
