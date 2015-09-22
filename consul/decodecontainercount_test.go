package consul

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecodeContainerCount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/kv/priority1-demand" {
			t.Fatalf("Url not as expected, have %s", r.URL.Path)
		}

		if r.Method != "GET" {
			t.Fatalf("Expected GET, have %s", r.Method)
		}

		h := w.Header()
		h.Set("Content-Type", "application/json")
		fmt.Fprintln(w, `[
	 	{
			 "CreateIndex": 8,
			 "ModifyIndex": 15,
			 "LockIndex": 0,
			 "Key": "priority1-demand",
			 "Flags": 0,
			 "Value": "OQ=="
	 	}
		]`)
	}))
	defer server.Close()

	d := NewDemandFromConsul()
	d.baseConsulUrl = server.URL

	count, err := d.GetDemand("priority1-demand")
	if err != nil {
		t.Fatalf("Error returned. %v", err)
	}
	if count != 9 {
		t.Fatalf("Decoded count not as expected. Have %d", count)
	}
}
