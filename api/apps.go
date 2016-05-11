package api

import (
	"encoding/json"

	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/metric"
	"github.com/microscaling/microscaling/target"
	"github.com/microscaling/microscaling/utils"
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
	MetricType        string          `json:"metricType"`
	Config            DockerAppConfig `json:"config"`
}

// TODO!! This is not really just Docker-specific as we have some target info in here too
type DockerAppConfig struct {
	Image           string `json:"image"`
	Command         string `json:"command"`
	PublishAllPorts bool   `json:"publishAllPorts"`
	QueueLength     int    `json:"targetQueueLength"`
	QueueName       string `json:"queueName"`
	TopicName       string `json:"topicName"`
	ChannelName     string `json:"channelName"`
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
			MaxDelta:      (a.MaxContainers - a.MinContainers),
			IsScalable:    true,

			// TODO!! Settings that need to be made configurable via the API.
			// Default PublishAllPorts to true.
			PublishAllPorts: true,
			// Set Network mode to host only. This won't work for load balancer metrics.
			NetworkMode: "host",
		}

		switch a.RuleType {
		case "Queue":
			task.Target = target.NewQueueLengthTarget(a.Config.QueueLength)
			switch a.MetricType {
			case "AzureQueue":
				task.Metric = metric.NewAzureQueueMetric(a.Config.QueueName)
			case "NSQ":
				task.Metric = metric.NewNSQMetric(a.Config.TopicName, a.Config.ChannelName)
			default:
				log.Errorf("Unexpected queue metricType %s", a.MetricType)
			}
		default:
			task.Target = target.NewRemainderTarget(a.MaxContainers)
			task.Metric = metric.NewNullMetric()
		}

		tasks = append(tasks, &task)
	}

	if err != nil {
		log.Debugf("Apps message: %v", appsMessage)
	}

	return
}

func GetApps(apiAddress string, userID string) (tasks []*demand.Task, maxContainers int, err error) {
	url := "http://" + apiAddress + "/apps/" + userID

	body, err := utils.GetJSON(url)
	if err != nil {
		log.Debugf("Failed to get /apps/: %v", err)
		return nil, 0, err
	}

	tasks, maxContainers, err = appsFromResponse(body)
	return
}
