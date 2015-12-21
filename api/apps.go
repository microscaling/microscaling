package api

import (
	"encoding/json"
	// "io"
	"log"

	"github.com/force12io/force12/demand"
)

type AppDescription struct {
	Name    string          `json:"name"`
	AppType string          `json:"type"`
	Config  DockerAppConfig `json:"config"`
}

type DockerAppConfig struct {
	Image           string `json:"image"`
	Command         string `json:"command"`
	PublishAllPorts bool   `json:"publish_all_ports"`
}

type dockerAppConfig DockerAppConfig

func (d *DockerAppConfig) UnmarshalJSON(b []byte) (err error) {
	c := dockerAppConfig{}
	err = json.Unmarshal(b, &c)
	if err == nil {
		*d = DockerAppConfig(c)
	}
	return
}

func appsFromResponse(b []byte) (tasks map[string]demand.Task, err error) {
	var apps []AppDescription
	err = json.Unmarshal(b, &apps)

	tasks = make(map[string]demand.Task)

	for _, a := range apps {
		name := a.Name
		task := demand.Task{
			Image:   a.Config.Image,
			Command: a.Config.Command,
			// TODO!! For now we will default turning on publishAllPorts, until we add this to the client-server interface
			PublishAllPorts: true,
		}

		tasks[name] = task
	}

	return
}

// Get /apps/ to receive a list of apps we'll be scaling
func GetApps(userID string) (tasks map[string]demand.Task, err error) {
	body, err := getJsonGet(userID, "/apps/")
	if err != nil {
		log.Printf("Failed to get /apps/: %v", err)
		return nil, err
	}

	tasks, err = appsFromResponse(body)
	return
}
