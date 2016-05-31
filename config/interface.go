package config

import (
	"github.com/microscaling/microscaling/demand"
	"github.com/op/go-logging"
)

// Config is an interface for retrieving task config - could be hardcoded, from the server, from a file etc
type Config interface {
	GetApps(userID string) (tasks []*demand.Task, maxContainers int, err error)
}

var log = logging.MustGetLogger("mssconfig")
