package api

import (
	"testing"
)

func TestGetBaseUrl(t *testing.T) {
	base := getBaseF12APIUrl()
	if base != "http://app.force12.io" || base != baseF12APIUrl {
		t.Fatalf("Maybe F12_METRICS_API_ADDRESS is set: %v | %v", base, baseF12APIUrl)
	}
}
