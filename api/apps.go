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
	Image   string `json:"image"`
	Command string `json:"command"`
}

type dockerAppConfig DockerAppConfig

// func Decode(r io.Reader) (x *AppDescription, err error) {
// 	x = new(AppDescription)
// 	err = json.NewDecoder(r).Decode(x)
// 	if err != nil {
// 		log.Printf("Error %v", err)
// 	}
// 	return
// }

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

	log.Println(string(body))
	tasks, err = appsFromResponse(body)
	return
}
