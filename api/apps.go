package api

import (
	"encoding/json"

	"github.com/microscaling/microscaling/demand"
)

type AppsMessage struct {
	UserID        string           `json:"name"`
	MaxContainers int              `json:"maxContainers"`
	Apps          []AppDescription `json:"apps"`
}

type AppDescription struct {
	Name              string          `json:"name"`
	Priority          int             `json:"priority"` // 1 is the highest, 0 means it's not scalable
	MinContainers     int             `json:"minContainers"`
	MaxContainers     int             `json:"maxContainers"`
	TargetQueueLength int             `json:"targetValue"`
	RuleType          string          `json:"ruleType"`
	AppType           string          `json:"appType"`
	Config            DockerAppConfig `json:"config"`
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

func appsFromResponse(b []byte) (tasks []*demand.Task, maxContainers int, err error) {
	var appsMessage AppsMessage

	err = json.Unmarshal(b, &appsMessage)
	if err != nil {
		log.Debugf("Error unmarshalling from %s", string(b[:]))
	}

	maxContainers = appsMessage.MaxContainers

	for _, a := range appsMessage.Apps {
		task := demand.Task{
			Name:          a.Name,
			Image:         a.Config.Image,
			Command:       a.Config.Command,
			Priority:      a.Priority,
			MinContainers: a.MinContainers,
			MaxContainers: a.MaxContainers,
			TargetType:    a.RuleType,
			// TODO!! For now we will default turning on publishAllPorts, until we add this to the client-server interface
			PublishAllPorts: true,
		}

		if a.RuleType == "Queue" {
			task.Target = a.TargetQueueLength
		}

		tasks = append(tasks, &task)
	}

	return
}

func GetApps(userID string) (tasks []*demand.Task, maxContainers int, err error) {
	body, err := getJsonGet(userID, "/v2/apps/")
	if err != nil {
		log.Debugf("Failed to get /v2/apps/: %v", err)
		return nil, 0, err
	}

	tasks, maxContainers, err = appsFromResponse(body)
	return
}
