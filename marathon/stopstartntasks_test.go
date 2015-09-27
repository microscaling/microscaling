package marathon

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStartStop(t *testing.T) {
	tests := []struct {
		expPayload   string
		statusCode   int
		expUrl       string
		app          string
		family       string
		demandcount  int
		currentcount int
		expErr       bool
	}{
		{
			app:          "myapp",
			family:       "bananas",
			demandcount:  99,
			currentcount: 37,
			expUrl:       "/myapp",
			expPayload: `{"instances":99}
`,
		},
		{
			app:          "myapp",
			family:       "bananas",
			demandcount:  99,
			currentcount: 37,
			expUrl:       "/myapp",
			statusCode:   500,
			expErr:       true,
			expPayload: `{"instances":99}
`,
		},
	}

	for _, test := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != test.expUrl {
				t.Fatalf("Expected root path, have %s", r.URL.Path)
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
			if string(data) != test.expPayload {
				t.Fatalf("post payload not as expected have %s, expected **%s**", string(data), test.expPayload)
			}
			w.WriteHeader(test.statusCode)
		}))
		defer server.Close()

		m := NewScheduler()
		m.baseMarathonUrl = server.URL

		log.Println("before start/stop: current, demand", test.currentcount, test.demandcount)
		err := m.StopStartNTasks(test.app, test.family, test.demandcount, &test.currentcount)
		log.Println("after start/stop: current, demand", test.currentcount, test.demandcount)

		if err != nil {
			if !test.expErr {
				t.Fatalf("Error. %v", err)
			}

			if test.currentcount == test.demandcount {
				t.Fatalf("Currentcount should not have been updated because we expected an error for this test")
			}
		} else if test.expErr {
			t.Fatalf("expected an error")
		} else {
			if test.currentcount != test.demandcount {
				t.Fatalf("Currentcount should have been updated")
			}

		}

	}
}
