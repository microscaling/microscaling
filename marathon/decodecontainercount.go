package marathon

import (
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
)

func (m *MarathonScheduler) GetContainerCount(key string) (int, error) {
	value, err := m.GetValuebyID(key)
	if err != nil {
		return 0, err
	}

	return m.DecodeContainerCount(value)
}

// DecodeContainerCount is called to decode the container count value retrieved
// as a string from Consul. It is a base64 encoded integer
func (m *MarathonScheduler) DecodeContainerCount(encoded string) (int, error) {
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
