package metric

import (
	"github.com/op/go-logging"
)

type Metric interface {
	UpdateCurrent()
	Current() int
}

var log = logging.MustGetLogger("mssmetric")
