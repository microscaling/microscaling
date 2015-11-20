package api

import (
	"log"
	"net/http/httptest"
	"reflect"
	"testing"

	"golang.org/x/net/websocket"
)

var tests = []struct {
	testJson  string
	expDemand []TaskDemand
}{
	{
		testJson: `{
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
		expDemand: []TaskDemand{
			{
				App:         "priority2",
				DemandCount: 3,
			},
			{
				App:         "priority1",
				DemandCount: 7,
			},
		},
	},
}

var testIndex int
var du chan []TaskDemand = make(chan []TaskDemand, 1)

func testServerDemand(ws *websocket.Conn) {
	err := Listen(ws, du)

	if err != nil {
		log.Printf("Error %v", err)
	}
}

func TestGetDemand(t *testing.T) {
	server := httptest.NewServer(websocket.Handler(testServerDemand))
	serverAddr = server.Listener.Addr().String()

	baseF12APIUrl = serverAddr
	ws, err := InitWebSocket()
	if err != nil {
		t.Fatal("dialing", err)
	}

	var result []TaskDemand

	for _, test := range tests {
		// Send message as if it were from the server
		_, err = ws.Write([]byte(test.testJson))

		// Listener function should send the received result here
		result = <-du
		expected := test.expDemand

		var rr map[string]int = make(map[string]int, 10)
		var ee map[string]int = make(map[string]int, 10)

		for _, v := range result {
			rr[v.App] = v.DemandCount
		}

		for _, v := range expected {
			ee[v.App] = v.DemandCount
		}

		if !reflect.DeepEqual(rr, ee) {
			log.Printf("Received %#v", result)
			log.Printf("Expected %#v", expected)
			t.Fatalf("Unexpected demand")
		}
	}

	ws.Close()
	server.Close()
}
