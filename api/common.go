// API between Force12 agent and server
package api

import (
	"fmt"
	"io/ioutil"
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
	Timeout: 15000 * time.Millisecond,
}

func getJsonGet(userID string, endpoint string) (body []byte, err error) {
	url := baseF12APIUrl + endpoint + userID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to build API GET request err %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to GET err %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Http error %d: %s", resp.StatusCode, resp.Status)
	}

	body, err = ioutil.ReadAll(resp.Body)

	return body, err
}
