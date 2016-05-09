// Package api describes the API between Microscaling agent and server
package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/op/go-logging"

	"golang.org/x/net/websocket"
)

var (
	log                 = logging.MustGetLogger("mssapi")
	baseAPIUrl          = GetBaseAPIUrl()
	debugTimeHTTPClient = (os.Getenv("MSS_TIME_HTTP_CLIENT") == "true")
	httpClient          = &http.Client{
		Timeout: 15000 * time.Millisecond,
	}
)

// GetBaseAPIUrl returns the server URL
func GetBaseAPIUrl() string {
	baseURL := os.Getenv("MSS_API_ADDRESS")
	if baseURL == "" {
		baseURL = "app.microscaling.com"
	}

	log.Infof("Sending results to %s", baseURL)
	return baseURL
}

// SetBaseAPIUrl only used for testing
func SetBaseAPIUrl(baseURL string) {
	baseAPIUrl = baseURL
}

func getJsonGet(userID string, endpoint string) (body []byte, err error) {
	url := "http://" + baseAPIUrl + endpoint + userID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to build API GET request err %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := timeHTTPClientDo(req)
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

func timeHTTPClientDo(req *http.Request) (resp *http.Response, err error) {
	var issuedAt time.Time
	if debugTimeHTTPClient {
		issuedAt = time.Now()
	}

	resp, err = httpClient.Do(req)

	if debugTimeHTTPClient {
		apiDuration := time.Since(issuedAt)
		log.Debugf("%v took %v", req.URL, apiDuration)
	}

	return
}

// InitWebSocket opens a websocket to the server
func InitWebSocket() (ws *websocket.Conn, err error) {
	origin := "http://localhost/"
	url := "ws://" + baseAPIUrl
	ws, err = websocket.Dial(url, "", origin)
	if err != nil {
		log.Errorf("Error getting the web socket: %v", err)
	}

	return ws, err
}
