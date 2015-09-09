package marathon

import (
	"os"
)

type MarathonScheduler struct {
	baseMarathonUrl string
	baseConsulUrl   string
}

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
