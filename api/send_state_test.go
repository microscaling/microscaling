package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
		expUrl   string
		expP1    int
		expP1Req int
		expP2    int
		expP2Req int
	}{
		{
			expUrl:   "/metrics/hello",
			expP1:    4,
			expP1Req: 3,
			expP2:    5,
			expP2Req: 7,
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

			payload := metricsPayload{}
			json.Unmarshal(data, &payload)

			for key, value := range payload.Tasks {
				log.Println(key, value)
				if value.App == "priority1" {
					if value.RunningCount != test.expP1 {
						t.Fatalf("Bad running for %s expected %d got %d", key, test.expP1, value.RunningCount)
					}
					if value.PendingCount != test.expP1Req {
						t.Fatalf("Bad pending for %s expected %d got %d", key, test.expP1Req, value.PendingCount)
					}
				} else if value.App == "priority2" {
					if value.RunningCount != test.expP2 {
						t.Fatalf("Bad running for %s expected %d got %d", key, test.expP2, value.RunningCount)
					}
					if value.PendingCount != test.expP2Req {
						t.Fatalf("Bad pending for %s expected %d got %d", key, test.expP2Req, value.PendingCount)
					}
				} else {
					t.Fatalf("Unexpected app name %s", value.App)
				}
			}
		}))
		defer server.Close()

		baseF12APIUrl = server.URL
		SendMetrics("hello", tasks)
		baseF12APIUrl = GetBaseF12APIUrl()
	}
}
