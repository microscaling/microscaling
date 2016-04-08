package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/op/go-logging"

	"github.com/microscaling/microscaling/api"
	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/scheduler"
	"github.com/microscaling/microscaling/scheduler/docker"
	"github.com/microscaling/microscaling/scheduler/toy"
)

type settings struct {
	schedulerType string
	sendMetrics   bool
	userID        string
	pullImages    bool
	dockerHost    string
}

func initLogging() {
	basicLogFormat := logging.MustStringFormatter(`%{color}%{level:.4s} %{time:15:04:05.000}: %{color:reset} %{message}`)
	detailLogFormat := logging.MustStringFormatter(`%{color}%{level:.4s} %{time:15:04:05.000} %{pid} %{shortfile}: %{color:reset} %{message}`)

	// "%{level:.1s}%{time:0102 15:04:05.999999} %{pid} %{shortfile}] %{message}"
	logComponents := getEnvOrDefault("MSS_LOG_DEBUG", "none")
	if strings.Contains(logComponents, "detail") {
		logging.SetFormatter(detailLogFormat)
	} else {
		logging.SetFormatter(basicLogFormat)
	}

	logBackend := logging.NewLogBackend(os.Stdout, "", 0)
	logging.SetBackend(logBackend)

	var components = []string{"mssengine", "mssagent", "mssapi", "mssdemand", "mssscheduler"}

	switch logComponents {
	case "all":
		for _, component := range components {
			logging.SetLevel(logging.DEBUG, component)
		}
	case "none":
		for _, component := range components {
			logging.SetLevel(logging.INFO, component)
		}
	default:
		for _, component := range components {
			if strings.Contains(logComponents, component) {
				logging.SetLevel(logging.DEBUG, component)
			} else {
				logging.SetLevel(logging.INFO, component)
			}
		}
	}
}

func getSettings() settings {
	var st settings
	st.schedulerType = getEnvOrDefault("MSS_SCHEDULER", "DOCKER")
	st.userID = getEnvOrDefault("MSS_USER_ID", "5k5gk")
	st.sendMetrics = (getEnvOrDefault("MSS_SEND_METRICS_TO_API", "true") == "true")
	st.pullImages = (getEnvOrDefault("MSS_PULL_IMAGES", "true") == "true")
	st.dockerHost = getEnvOrDefault("DOCKER_HOST", "unix:///var/run/docker.sock")
	return st
}

func getScheduler(st settings) (scheduler.Scheduler, error) {
	var s scheduler.Scheduler

	switch st.schedulerType {
	case "DOCKER":
		log.Info("Scheduling with Docker remote API")
		s = docker.NewScheduler(st.pullImages, st.dockerHost)
	case "ECS":
		return nil, fmt.Errorf("Scheduling with ECS not yet supported. Tweet with hashtag #MicroscaleECS if you'd like us to add this next!")
	case "KUBERNETES":
		return nil, fmt.Errorf("Scheduling with Kubernetes not yet supported. Tweet with hashtag #MicroscaleK8S if you'd like us to add this next!")
	case "MESOS":
		return nil, fmt.Errorf("Scheduling with Mesos / Marathon not yet supported. Tweet with hashtag #MicroscaleMesos if you'd like us to add this next!")
	case "NOMAD":
		return nil, fmt.Errorf("Scheduling with Nomad not yet supported. Tweet with hashtag #MicroscaleNomad if you'd like us to add this next!")
	case "TOY":
		log.Info("Scheduling with toy scheduler")
		s = toy.NewScheduler()
	default:
		return nil, fmt.Errorf("Bad value for MSS_SCHEDULER: %s", st.schedulerType)
	}

	if s == nil {
		return nil, fmt.Errorf("No scheduler")
	}

	return s, nil
}

func getTasks(st settings) map[string]demand.Task {
	var t map[string]demand.Task

	// Get the tasks that have been configured by this user
	t, err := api.GetApps(st.userID)
	if err != nil {
		log.Errorf("Error getting tasks: %v", err)
	}

	log.Debug(t)
	return t
}

func getEnvOrDefault(name string, defaultValue string) string {
	v := os.Getenv(name)
	if v == "" {
		v = defaultValue
	}

	return v
}
