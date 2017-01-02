package config

import (
	"strconv"
	"strings"

	microbadger "github.com/microscaling/microbadger/api"

	"github.com/microscaling/microscaling/api"
	"github.com/microscaling/microscaling/demand"
)

// LabelConfig is used when we retrieve get image config from the Microscaling server and then
// get label config from MicroBadger APIs
type LabelConfig struct {
	APIAddress string
}

// compile-time assert that we implement the right interface
var _ Config = (*LabelConfig)(nil)

// NewLabelConfig gets a new LabelConfig
func NewLabelConfig(APIAddress string) *LabelConfig {
	return &LabelConfig{
		APIAddress: APIAddress,
	}
}

// GetApps retrieves task config from the server using the API, and then gets scaling parameters from labels using MicroBadger
func (l *LabelConfig) GetApps(userID string) (tasks []*demand.Task, maxContainers int, err error) {
	tasks, maxContainers, err = api.GetApps(l.APIAddress, userID)
	for _, task := range tasks {
		labels, err := microbadger.GetLabels(task.Image)
		if err != nil {
			log.Errorf("Failed to get labels for %s: %v", task.Image, err)
		} else {
			parseLabels(task, labels)
		}
	}
	return
}

func parseLabels(task *demand.Task, labels map[string]string) {
	// Make sure there's a lower-case version of all labels (don't overwrite a
	// lower-case one if it's already there)
	for k, v := range labels {
		kl := strings.ToLower(k)
		if kl != k {
			if _, ok := labels[kl]; !ok {
				labels[kl] = v
			}
		}
	}

	if isScalable, ok := labels["com.microscaling.is-scalable"]; ok {
		if b, err := strconv.ParseBool(isScalable); err == nil {
			task.IsScalable = b
		}
	}

	v, err := parseIntLabel(labels, "com.microscaling.priority")
	if err == nil {
		task.Priority = v
	}

	v, err = parseIntLabel(labels, "com.microscaling.max-delta")
	if err == nil {
		task.MaxDelta = v
	}

	v, err = parseIntLabel(labels, "com.microscaling.min-containers")
	if err == nil {
		task.MinContainers = v
	}

	v, err = parseIntLabel(labels, "com.microscaling.max-containers")
	if err == nil {
		task.MaxContainers = v
	}
}

func parseIntLabel(labels map[string]string, key string) (intVal int, err error) {
	if val, ok := labels[key]; ok {
		intVal, err = strconv.Atoi(val)
	}

	if err != nil {
		log.Infof("Ignoring bad value for label %s", key)
	}
	return
}
