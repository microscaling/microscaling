package metric

import (
	"github.com/op/go-logging"
)

// Metric is something we measure. Each task is associated with a Metric and a Target that we want the Metric to stay close to.
type Metric interface {
	UpdateCurrent()
	Current() int
}

var log = logging.MustGetLogger("mssmetric")
