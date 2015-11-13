package api

import (
	"testing"
)

func TestGetBaseUrl(t *testing.T) {
	base := GetBaseF12APIUrl()
	if base != "app.force12.io" || base != baseF12APIUrl {
		t.Fatalf("Maybe F12_API_ADDRESS is set: %v | %v", base, baseF12APIUrl)
	}
}
