package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/force12io/force12/demand"
)

func TestSendState(t *testing.T) {
	var tasks map[string]demand.Task = make(map[string]demand.Task)

	tasks["priority1"] = demand.Task{Demand: 8, Requested: 3, Running: 4}
	tasks["priority2"] = demand.Task{Demand: 2, Requested: 7, Running: 5}

	tests := []struct {
		expUrl           string
		expMaxContainers int
		expP1            int
		expP2            int
	}{
		{
			expUrl:           "/metrics/hello",
			expMaxContainers: 99,
			expP1:            4,
			expP2:            5,
		},
	}

	for _, test := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != test.expUrl {
				t.Fatalf("Expected %s, have %s", test.expUrl, r.URL.Path)
			}

			if r.Method != "PUT" {
				t.Fatalf("expected PUT, have %s", r.Method)
			}

			if r.Header.Get("Content-Type") != "application/json" {
				t.Fatalf("Content type not as expected, have %s", r.Header.Get("Content-Type"))
			}
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Failed to read post body. %v", err)
			}
			payload := sendStatePayload{}
			json.Unmarshal(data, &payload)

			if payload.MaxContainers != test.expMaxContainers {
				t.Fatalf("Wrong max container count %d", payload.MaxContainers)
			}
			if payload.Priority1Running != test.expP1 {
				t.Fatalf("Wrong P1 count %d", payload.Priority1Running)
			}
			if payload.Priority2Running != test.expP2 {
				t.Fatalf("Wrong P2 count %d", payload.Priority2Running)
			}
		}))
		defer server.Close()

		baseF12APIUrl = server.URL
		SendState("hello", tasks, test.expMaxContainers)
		baseF12APIUrl = GetBaseF12APIUrl()
	}
}
