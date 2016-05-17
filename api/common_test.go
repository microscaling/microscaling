package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/utils"
	"golang.org/x/net/websocket"
)

var serverAddr string

func testServer(ws *websocket.Conn) {
	log.Debugf("Received something")
}

func TestInitWebSocket(t *testing.T) {
	server := httptest.NewServer(websocket.Handler(testServer))
	serverAddr = server.Listener.Addr().String()

	ws, err := utils.InitWebSocket(serverAddr)
	if err != nil {
		t.Fatal("dialing", err)
	}

	msg := []byte("hello, world\n")
	if _, err := ws.Write(msg); err != nil {
		t.Errorf("Write: %v", err)
	}

	ws.Close()
	server.Close()
}

// Utility for checking GET requests
func DoTestGetJSON(t *testing.T, expURL string, success bool, testJSON string) (server *httptest.Server) {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expURL {
			t.Fatalf("Expected %s, have %s", expURL, r.URL.Path)
		}

		if r.Method != "GET" {
			t.Fatalf("expected GET, have %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Content type not as expected, have %s", r.Header.Get("Content-Type"))
		}

		if success {
			w.Write([]byte(testJSON))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))

	return server
}

// Utility for checking that tasks are updated to be what we expect
func CheckReturnedTasks(t *testing.T, tasks map[string]demand.Task, returnedTasks []*demand.Task) {
	for _, rt := range returnedTasks {
		tt, ok := tasks[rt.Name]
		if !ok {
			t.Fatalf("Unexpected app name %v", rt.Name)
		}

		if tt.Image != rt.Image {
			t.Fatalf("Image: expected %s got %s", tt.Image, rt.Image)
		}
		if tt.Command != rt.Command {
			t.Fatalf("Command: expected %s got %s", tt.Command, rt.Command)
		}
		if tt.Demand != rt.Demand {
			t.Fatalf("Demand: expected %s got %s", tt.Demand, rt.Demand)
		}
		if tt.Requested != rt.Requested {
			t.Fatalf("Requested: expected %s got %s", tt.Requested, rt.Requested)
		}
		if tt.Running != rt.Running {
			t.Fatalf("Requested: expected %s got %s", tt.Requested, rt.Requested)
		}
	}
}
