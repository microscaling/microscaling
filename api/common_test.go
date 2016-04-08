package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/microscaling/microscaling/demand"
	"golang.org/x/net/websocket"
)

var serverAddr string

func TestGetBaseUrl(t *testing.T) {
	base := GetBaseAPIUrl()
	if base != "app.microscaling.com" || base != baseAPIUrl {
		t.Fatalf("Maybe MSS_API_ADDRESS is set: %v | %v", base, baseAPIUrl)
	}
}

func testServer(ws *websocket.Conn) {
	log.Debugf("Received something")
}

func TestInitWebSocket(t *testing.T) {
	server := httptest.NewServer(websocket.Handler(testServer))
	serverAddr = server.Listener.Addr().String()

	baseAPIUrl = serverAddr
	ws, err := InitWebSocket()
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
func DoTestGetJson(t *testing.T, expUrl string, success bool, testJson string) (server *httptest.Server) {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expUrl {
			t.Fatalf("Expected %s, have %s", expUrl, r.URL.Path)
		}

		if r.Method != "GET" {
			t.Fatalf("expected GET, have %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Content type not as expected, have %s", r.Header.Get("Content-Type"))
		}

		if success {
			w.Write([]byte(testJson))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))

	return server
}

// Utility for checking that tasks are updated to be what we expect
func CheckReturnedTasks(t *testing.T, tasks map[string]demand.Task, returned_tasks map[string]demand.Task) {
	for name, rt := range returned_tasks {
		tt, ok := tasks[name]
		if !ok {
			t.Fatalf("Unexpected app name %v", name)
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
