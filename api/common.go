// API between Microscaling agent and server
package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

func GetBaseAPIUrl() string {
	baseUrl := os.Getenv("MSS_API_ADDRESS")
	if baseUrl == "" {
		baseUrl = "app.microscaling.com"
	}

	log.Printf("Sending results to %s", baseUrl)
	return baseUrl
}

// SetBaseAPIUrl only used for testing
func SetBaseAPIUrl(baseurl string) {
	baseAPIUrl = baseurl
}

var baseAPIUrl string = GetBaseAPIUrl()
var httpClient *http.Client = &http.Client{
	Timeout: 15000 * time.Millisecond,
}
var debugTimeHttpClient bool = (os.Getenv("MSS_TIME_HTTP_CLIENT") == "true")

func getJsonGet(userID string, endpoint string) (body []byte, err error) {
	url := "http://" + baseAPIUrl + endpoint + userID

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

func InitWebSocket() (ws *websocket.Conn, err error) {
	origin := "http://localhost/"
	url := "ws://" + baseAPIUrl
	ws, err = websocket.Dial(url, "", origin)
	if err != nil {
		log.Printf("Error getting the web socket: %v", err)
	}

	return ws, err
}
