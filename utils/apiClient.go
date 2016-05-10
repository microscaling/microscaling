// Package utils contains common shared code.
package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/op/go-logging"
	"golang.org/x/net/websocket"
)

var log = logging.MustGetLogger("mssmetric")

var (
	httpClient = &http.Client{
		// TODO Make timeout configurable.
		Timeout: 10 * time.Second,
	}
)

// GetJSON makes a GET request to a REST API and returns the JSON response.
func GetJSON(url string) (body []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Failed to build API GET request err %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Errorf("API request failed %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("GET request failed %s %d: %s", url, resp.StatusCode, resp.Status)
	}

	body, err = ioutil.ReadAll(resp.Body)

	return body, err
}

// PutJSON makes a PUT request to a REST API and submits the JSON payload.
func PutJSON(url string, payload *bytes.Buffer) (status int, err error) {
	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		log.Errorf("Failed to build PUT request err %v", err)
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("API request failed %v", url, resp.StatusCode, err)
		return -1, err
	} else {
		resp.Body.Close()
		return resp.StatusCode, err
	}
}

// InitWebSocket opens a websocket connection to the provided address.
func InitWebSocket(apiAddress string) (ws *websocket.Conn, err error) {
	// TODO Set actual origin.
	origin := "http://localhost/"
	url := "ws://" + apiAddress
	ws, err = websocket.Dial(url, "", origin)
	if err != nil {
		log.Errorf("Error getting the web socket: %v", err)
	}

	return ws, err
}
