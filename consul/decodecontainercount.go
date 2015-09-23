package consul

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"

	"bitbucket.org/force12io/force12-scheduler/demand"
)

type DemandFromConsul struct {
	baseConsulUrl string
}

// check that we implement the demand interface
var _ demand.Input = (*DemandFromConsul)(nil)

func NewDemandModel() *DemandFromConsul {
	return &DemandFromConsul{
		baseConsulUrl: getBaseConsulUrl(),
	}
}

func getBaseConsulUrl() string {
	baseUrl := os.Getenv("CONSUL_ADDRESS")
	if baseUrl == "" {
		baseUrl = "http://marathon.force12.io:8500"
	}
	return baseUrl
}

func (d *DemandFromConsul) GetDemand(key string) (int, error) {
	encoded, err := d.GetValuebyID(key)
	if err != nil {
		return 0, err
	}

	// Container count in consul is a base64 encoded integer
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return 0, fmt.Errorf("Failed to base64decode container count. %v", err)
	}

	containerCount, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, fmt.Errorf("Failed to convert container count. %v", err)
	}

	log.Printf("Container count: %d", containerCount)
	return containerCount, nil
}
