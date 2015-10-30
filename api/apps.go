// API between Force12 agent and server
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/force12io/force12/demand"
)

type AppDescription struct {
	Name    string          `json:"name"`
	AppType string          `json:"type"`
	Config  DockerAppConfig `json:"config"`
}

type DockerAppConfig struct {
	Image   string `json:"image"`
	Command string `json:"command"`
}

type dockerAppConfig DockerAppConfig

func Decode(r io.Reader) (x *AppDescription, err error) {
	x = new(AppDescription)
	err = json.NewDecoder(r).Decode(x)
	if err != nil {
		log.Printf("Error %v", err)
	}
	return
}

func (d *DockerAppConfig) UnmarshalJSON(b []byte) (err error) {
	c := dockerAppConfig{}
	err = json.Unmarshal(b, &c)
	if err == nil {
		*d = DockerAppConfig(c)
	}
	return
}

func appsFromResponse(b []byte) (tasks map[string]demand.Task, err error) {
	var a []AppDescription
	err = json.Unmarshal(b, &a)

	tasks = make(map[string]demand.Task)

	for _, a := range a {
		name := a.Name
		task := demand.Task{
			Image: a.Config.Image,
		}

		tasks[name] = task
	}

	return
}

func GetApps(userID string) (tasks map[string]demand.Task, err error) {
	url := baseF12APIUrl + "/apps/" + userID

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
	body, err := ioutil.ReadAll(resp.Body)

	tasks, err = appsFromResponse(body)
	return
}
