package config

import (
	"github.com/microscaling/microscaling/api"
	"github.com/microscaling/microscaling/demand"
)

// EnvVarConfig is used when we retrieve config from an environment variable.
type EnvVarConfig struct {
	EnvVarValue string
}

// compile-time assert that we implement the right interface.
var _ Config = (*EnvVarConfig)(nil)

// NewEnvVarConfig gets a new EnvVarConfig.
func NewEnvVarConfig(EnvVarValue string) *EnvVarConfig {
	return &EnvVarConfig{
		EnvVarValue: EnvVarValue,
	}
}

// GetApps retrieves tasks from the environment variable.
func (e *EnvVarConfig) GetApps(userID string) (tasks []*demand.Task, maxContainers int, err error) {
	data := ([]byte(e.EnvVarValue))

	return api.AppsFromData(data)
}
