package marathon

import (
	"os"

	"bitbucket.org/force12io/force12-scheduler/scheduler"
)

type MarathonScheduler struct {
	baseMarathonUrl string
	baseConsulUrl   string
}

// compile-time assert that we implement the right interface
var _ scheduler.Scheduler = (*MarathonScheduler)(nil)

func NewMarathonScheduler() *MarathonScheduler {
	return &MarathonScheduler{
		baseMarathonUrl: getBaseMarathonUrl(),
		baseConsulUrl:   getBaseConsulUrl(),
	}
}

func getBaseConsulUrl() string {
	baseUrl := os.Getenv("CONSUL_ADDRESS")
	if baseUrl == "" {
		baseUrl = "http://marathon.force12.io:8500"
	}
	return baseUrl
}

func getBaseMarathonUrl() string {
	baseUrl := os.Getenv("MARATHON_ADDRESS")
	if baseUrl == "" {
		baseUrl = "http://marathon.force12.io:8080"
	}
	return baseUrl + "/v2/apps"
}
