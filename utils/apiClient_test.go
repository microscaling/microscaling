package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/websocket"
)

func testServer(ws *websocket.Conn) {
	log.Debugf("Received something")
}

func TestInitWebSocket(t *testing.T) {
	server := httptest.NewServer(websocket.Handler(testServer))
	serverAddr := server.Listener.Addr().String()

	ws, err := InitWebSocket(serverAddr)
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

func DoTestJSON(t *testing.T, method string, testJSON string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Fatalf("expected method, have %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Content type not as expected, have %s", r.Header.Get("Content-Type"))
		}

		if method == "GET" {
			w.Write([]byte(testJSON))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))

	return server
}

func TestGetJSON(t *testing.T) {
	tests := []struct {
		testJSON string
	}{
		{
			testJSON: "{'hello':'world'}",
		}, {
			testJSON: "",
		},
	}

	for _, test := range tests {
		server := DoTestJSON(t, "GET", test.testJSON)
		defer server.Close()

		_, err := GetJSON(server.URL)
		if err != nil {
			t.Fatalf("Failed unexpectedly: %v", err.Error())
		}
	}

	for _, test := range tests {
		server := DoTestJSON(t, "PUT", test.testJSON)
		defer server.Close()

		_, err := PutJSON(server.URL, bytes.NewBufferString(test.testJSON))
		if err != nil {
			t.Fatalf("PUT Failed unexpectedly: %v", err.Error())
		}

	}

}
