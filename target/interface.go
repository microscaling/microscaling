package target

import (
	"github.com/op/go-logging"
)

// Target is what we want the Metric to match. Each task has a Metric and a Target.
type Target interface {
	Meeting(int) bool
	Exceeding(int) bool
	Delta(int) int
}

var log = logging.MustGetLogger("msstarget")
