package main

import (
	"fmt"
	"log"
	"os"

	"github.com/microscaling/microscaling/api"
	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/docker"
	"github.com/microscaling/microscaling/scheduler"
	"github.com/microscaling/microscaling/toy_scheduler"
)

type settings struct {
	schedulerType string
	sendMetrics   bool
	userID        string
	pullImages    bool
	dockerHost    string
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
		log.Println("Scheduling with Docker remote API")
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
		log.Println("Scheduling with toy scheduler")
		s = toy_scheduler.NewScheduler()
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
		log.Printf("Error getting tasks: %v", err)
	}

	log.Println(t)
	return t
}

func getEnvOrDefault(name string, defaultValue string) string {
	v := os.Getenv(name)
	if v == "" {
		v = defaultValue
	}

	return v
}
