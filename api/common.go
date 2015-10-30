package api

import (
	"log"
	"net/http"
	"os"
	"time"
)

func getBaseF12APIUrl() string {
	baseUrl := os.Getenv("F12_METRICS_API_ADDRESS")
	if baseUrl == "" {
		baseUrl = "http://app.force12.io"
	}

	log.Printf("Sending results to %s", baseUrl)
	return baseUrl
}

var baseF12APIUrl string = getBaseF12APIUrl()
var httpClient *http.Client = &http.Client{
	Timeout: 30000 * time.Millisecond,
}
