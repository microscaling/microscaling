package config

import (
	"github.com/microscaling/microscaling/api"
	"github.com/microscaling/microscaling/demand"
)

// ServerConfig is used when we retrieve config over the API from the server
type ServerConfig struct{}

// compile-time assert that we implement the right interface
var _ Config = (*ServerConfig)(nil)

// NewServerConfig gets a new ServerConfig
func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

// GetApps retrieves task config from the server using the API
func (s *ServerConfig) GetApps(userID string) (tasks []*demand.Task, maxContainers int, err error) {
	tasks, maxContainers, err = api.GetApps(userID)
	return
}
