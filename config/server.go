package config

import (
	"github.com/microscaling/microscaling/api"
	"github.com/microscaling/microscaling/demand"
)

// ServerConfig is used when we retrieve config over the API from the server
type ServerConfig struct {
	APIAddress string
}

// compile-time assert that we implement the right interface
var _ Config = (*ServerConfig)(nil)

// NewServerConfig gets a new ServerConfig
func NewServerConfig(APIAddress string) *ServerConfig {
	return &ServerConfig{
		APIAddress: APIAddress,
	}
}

// GetApps retrieves task config from the server using the API
func (s *ServerConfig) GetApps(userID string) (tasks []*demand.Task, maxContainers int, err error) {
	tasks, maxContainers, err = api.GetApps(s.APIAddress, userID)
	return
}
