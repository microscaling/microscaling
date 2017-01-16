package config

import (
	"fmt"
	"strconv"
	"strings"

	microbadger "github.com/microscaling/microbadger/api"

	"github.com/microscaling/microscaling/api"
	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/utils"
)

// LabelConfig is used when we retrieve get image config from the Microscaling server and then
// get label config from MicroBadger APIs
type LabelConfig struct {
	APIAddress    string
	KubeConfig    string
	KubeNamespace string
}

// compile-time assert that we implement the right interface
var _ Config = (*LabelConfig)(nil)

// NewLabelConfig gets a new LabelConfig
func NewLabelConfig(APIAddress string) *LabelConfig {
	return &LabelConfig{
		APIAddress: APIAddress,
	}
}

// NewKubeLabelConfig gets a new LabelConfig for Kubernetes
func NewKubeLabelConfig(APIAddress string, KubeConfig string, KubeNamespace string) *LabelConfig {
	return &LabelConfig{
		APIAddress:    APIAddress,
		KubeConfig:    KubeConfig,
		KubeNamespace: KubeNamespace,
	}
}

// GetApps retrieves task config from the server using the API, and then gets scaling parameters from labels using MicroBadger
func (l *LabelConfig) GetApps(userID string) (tasks []*demand.Task, maxContainers int, err error) {
	tasks, maxContainers, err = api.GetApps(l.APIAddress, userID)
	for _, task := range tasks {
		if l.KubeNamespace != "" {
			task.Image, err = l.getImageFromKubeDeployment(task.Name)
			if err != nil {
				log.Errorf("Failed to get image for deployment %s: %v", task.Image, err)
				return
			}
		}

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

func (l *LabelConfig) getImageFromKubeDeployment(appName string) (imageName string, err error) {
	clientset, err := utils.NewKubeClientset(l.KubeConfig, l.KubeNamespace)
	if err != nil {
		log.Errorf("Error creating Kubernetes clientset: %v", err)
		return
	}

	d, err := clientset.Extensions().Deployments(l.KubeNamespace).Get(appName)
	if err != nil {
		log.Errorf("Error getting deployment %s: %v", appName, err)
		return
	}

	podSpec := d.Spec.Template.Spec
	containers := len(podSpec.Containers)

	if containers == 1 {
		imageName = podSpec.Containers[0].Image
		log.Debugf("Got image %s for deployment %s", imageName, appName)
	} else {
		// TODO!! Support pods with multiple containers
		return "", fmt.Errorf("Error expected 1 container per pod but found %d", containers)
	}

	return imageName, err
}
