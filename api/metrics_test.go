package api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/microscaling/microscaling/demand"
	"golang.org/x/net/websocket"
)

var global_t *testing.T

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

var mtests = []struct {
	expMetrics metricsPayload
}{
	{expMetrics: metricsPayload{
		User: "hello",
		Metrics: metrics{
			Tasks: []taskMetrics{
				{App: "priority1",
					RunningCount: 4,
					PendingCount: 3,
				},
				{App: "priority2",
					RunningCount: 5,
					PendingCount: 7},
			},
		},
	},
	},
}

var currentTest int = 0

func testServerMetrics(ws *websocket.Conn) {
	var b []byte
	b = make([]byte, 1000)
	length, _ := ws.Read(b)

	var m metricsPayload
	_ = json.Unmarshal(b[:length], &m)

	test := mtests[currentTest]

	if m.User != test.expMetrics.User {
		global_t.Fatalf("Unexpected user")
	}

	for _, v := range m.Metrics.Tasks {
		appFound := false
		for _, vv := range test.expMetrics.Metrics.Tasks {
			if vv.App == v.App {
				appFound = true
				if v.PendingCount != vv.PendingCount || v.RunningCount != vv.RunningCount {
					log.Debugf("%#v", test.expMetrics.Metrics.Tasks)
					log.Debugf("%#v", m.Metrics.Tasks)
					global_t.Fatalf("Unexpected values")
				}
			}
		}

		if !appFound {
			global_t.Fatalf("Received unexpected metric for %s", v.App)
		}
	}
}

func TestSendMetrics(t *testing.T) {
	var tasks demand.Tasks
	tasks.Tasks = make([]*demand.Task, 2)

	tasks.Tasks[0] = &demand.Task{Name: "priority1", Demand: 8, Requested: 3, Running: 4}
	tasks.Tasks[1] = &demand.Task{Name: "priority2", Demand: 2, Requested: 7, Running: 5}

	global_t = t

	for testIndex, _ = range tests {
		server := httptest.NewServer(websocket.Handler(testServerMetrics))
		serverAddr = server.Listener.Addr().String()

		baseAPIUrl = serverAddr
		ws, err := InitWebSocket()
		if err != nil {
			t.Fatal("dialing", err)
		}

		SendMetrics(ws, "hello", &tasks)

		ws.Close()
		server.Close()
	}
}
