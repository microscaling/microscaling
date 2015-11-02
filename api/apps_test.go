package api

import (
	"encoding/json"
	// "io/ioutil"
	// "bytes"
	// "log"
	// "net/http"
	// "net/http/httptest"
	"testing"

	"github.com/force12io/force12/demand"
)

func TestGetAppsDecode(t *testing.T) {
	var bb = []byte(`{"image": "my image", "command": "do it"}`)
	var d = DockerAppConfig{}
	_ = json.Unmarshal(bb, &d)
	if d.Image != "my image" {
		t.Fatalf("Didn't decode image")
	}

	if d.Command != "do it" {
		t.Fatalf("Didn't decode command")
	}

	var response string = `[{"name":"priority1","type":"Docker","config":{"image":"force12io/priority-1:latest","command":"/run.sh"}},{"name":"priority2","type":"Docker","config":{"image":"force12io/priority-2:latest","command":"/run.sh"}}]`
	var b = []byte(response)

	var a []AppDescription
	_ = json.Unmarshal(b, &a)

	var apps map[string]demand.Task
	apps, _ = appsFromResponse(b)

	p1 := apps["priority1"]
	if p1.Image != "force12io/priority-1:latest" {
		t.Fatalf("Bad image")
	}
	p2 := apps["priority2"]
	if p2.Image != "force12io/priority-2:latest" {
		t.Fatalf("Bad image")
	}
}
