// Package api defines API between Microscaling agent and server
package api

import (
	"github.com/op/go-logging"
)

type metricsPayload struct {
	User    string  `json:"user"`
	Metrics metrics `json:"metrics"`
}

type metrics struct {
	Tasks     []taskMetrics `json:"tasks"`
	CreatedAt int64         `json:"createdAt"`
}

type taskMetrics struct {
	App          string `json:"app"`
	RunningCount int    `json:"runningCount"`
	PendingCount int    `json:"pendingCount"`
	Metric       int    `json:"metric,omitempty"`
}

var log = logging.MustGetLogger("mssapi")
