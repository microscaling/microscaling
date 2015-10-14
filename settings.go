package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/force12io/force12/demand"
	"github.com/force12io/force12/docker"
	"github.com/force12io/force12/rng"
	"github.com/force12io/force12/scheduler"
	"github.com/force12io/force12/toy_scheduler"
)

type settings struct {
	demandModelType string
	schedulerType   string
	sendstate       bool
	userID          string
	demandInterval  time.Duration
	demandDelta     int
	maxContainers   int
}

func get_settings() settings {
	var st settings
	st.demandModelType = getEnvOrDefault("F12_DEMAND_MODEL", "RNG")
	st.schedulerType = getEnvOrDefault("F12_SCHEDULER", "DOCKER")
	st.userID = getEnvOrDefault("F12_USER_ID", "5k5gk")
	st.sendstate = (getEnvOrDefault("F12_SEND_STATE_TO_API", "true") == "true")
	st.demandDelta, _ = strconv.Atoi(getEnvOrDefault("F12_DEMAND_DELTA", "3"))
	st.maxContainers, _ = strconv.Atoi(getEnvOrDefault("F12_MAXIMUM_CONTAINERS", "9"))
	demandIntervalMs, _ := strconv.Atoi(getEnvOrDefault("F12_DEMAND_CHANGE_INTERVAL_MS", "3000"))
	st.demandInterval = time.Duration(demandIntervalMs) * time.Millisecond
	return st
}

func get_demand_input(st settings) (demand.Input, error) {
	var di demand.Input

	switch st.demandModelType {
	case "CONSUL":
		return nil, fmt.Errorf("Demand metric from Consul not yet supported")
	case "RNG":
		log.Println("Random demand generation")
		di = rng.NewDemandModel(st.demandDelta, st.maxContainers)
	default:
		return nil, fmt.Errorf("Bad value for F12_DEMAND_MODEL: %s", st.demandModelType)
	}

	return di, nil
}

func get_scheduler(st settings) (scheduler.Scheduler, error) {
	var s scheduler.Scheduler

	switch st.schedulerType {
	case "DOCKER":
		log.Println("Scheduling with Docker remote API")
		s = docker.NewScheduler()
	case "ECS":
		return nil, fmt.Errorf("Scheduling with ECS not yet supported. Tweet with hashtag #F12ECS if you'd like us to add this next!")
	case "KUBERNETES":
		return nil, fmt.Errorf("Scheduling with Kubernetes not yet supported. Tweet with hashtag #F12Kubernetes if you'd like us to add this next!")
	case "MESOS":
		return nil, fmt.Errorf("Scheduling with Mesos / Marathon not yet supported. Tweet with hashtag #F12Mesos if you'd like us to add this next!")
	case "NOMAD":
		return nil, fmt.Errorf("Scheduling with Nomad not yet supported. Tweet with hashtag #F12Nomad if you'd like us to add this next!")
	case "TOY":
		log.Println("Scheduling with toy scheduler")
		s = toy_scheduler.NewScheduler()
	default:
		return nil, fmt.Errorf("Bad value for F12_SCHEDULER: %s", st.schedulerType)
	}

	return s, nil
}

func get_tasks(st settings) map[string]demand.Task {
	var t map[string]demand.Task

	p1TaskName = getEnvOrDefault("F12_PRIORITY1_TASK", p1TaskName)
	p2TaskName = getEnvOrDefault("F12_PRIORITY2_TASK", p2TaskName)
	p1FamilyName = getEnvOrDefault("F12_PRIORITY1_FAMILY", p1FamilyName)
	p2FamilyName = getEnvOrDefault("F12_PRIORITY2_FAMILY", p2FamilyName)
	p1Image = getEnvOrDefault("F12_PRIORITY1_IMAGE", p1Image)
	p2Image = getEnvOrDefault("F12_PRIORITY2_IMAGE", p2Image)

	t = make(map[string]demand.Task)

	t[p1TaskName] = demand.Task{
		FamilyName: p1FamilyName,
		Demand:     st.maxContainers / 2,
		Requested:  0,
		Image:      p1Image,
	}

	t[p2TaskName] = demand.Task{
		FamilyName: p2FamilyName,
		Demand:     st.maxContainers - (st.maxContainers / 2),
		Requested:  0,
		Image:      p2Image,
	}

	return t
}

func getEnvOrDefault(name string, defaultValue string) string {
	v := os.Getenv(name)
	if v == "" {
		v = defaultValue
	}

	return v
}
