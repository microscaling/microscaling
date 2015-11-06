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

func GetBaseF12APIUrl() string {
	baseUrl := os.Getenv("F12_API_ADDRESS")
	if baseUrl == "" {
		baseUrl = "http://app.force12.io"
	}

	log.Printf("Sending results to %s", baseUrl)
	return baseUrl
}

// SetBaseF12APIUrl only used for testing
func SetBaseF12APIUrl(baseurl string) {
	baseF12APIUrl = baseurl
}

var baseF12APIUrl string = GetBaseF12APIUrl()
var httpClient *http.Client = &http.Client{
	Timeout: 15000 * time.Millisecond,
}
var debugTimeHttpClient bool = (os.Getenv("F12_TIME_HTTP_CLIENT") == "true")

func getJsonGet(userID string, endpoint string) (body []byte, err error) {
	url := baseF12APIUrl + endpoint + userID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to build API GET request err %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := timeHttpClientDo(req)
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

func timeHttpClientDo(req *http.Request) (resp *http.Response, err error) {
	var issuedAt time.Time
	if debugTimeHttpClient {
		issuedAt = time.Now()
	}

	resp, err = httpClient.Do(req)

	if debugTimeHttpClient {
		apiDuration := time.Since(issuedAt)
		log.Printf("%v took %v", req.URL, apiDuration)
	}

	return
}
