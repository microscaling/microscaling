package config

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServerConfig(t *testing.T) {
	expURL := "/apps/world"
	testJSON := `{
			"name": "world",
			"maxContainers": 10,
			"apps": [
			      {
			          "name": "priority1",
			          "appType": "Docker",
			          "config": {
			              "image": "firstimage"
			          }
			      },
			      {
			          "name": "priority2",
			          "appType": "Docker",
			          "config": {
			              "image": "anotherimage",
			              "command": "do this"
			          }
			      }
			]}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expURL {
			t.Fatalf("Expected %s, have %s", expURL, r.URL.Path)
		}

		if r.Method != "GET" {
			t.Fatalf("expected GET, have %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Content type not as expected, have %s", r.Header.Get("Content-Type"))
		}

		w.Write([]byte(testJSON))

	}))
	defer server.Close()

	baseAPIURL := strings.Replace(server.URL, "http://", "", 1)

	c := NewServerConfig(baseAPIURL)
	tasks, maxC, err := c.GetApps("world")

	if len(tasks) != 2 {
		t.Fatal("Expected two tasks")
	}

	if maxC != 10 {
		t.Fatal("Expected max containers 10")
	}

	if err != nil {
		t.Fatal("Should succeed")
	}
}
