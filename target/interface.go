package target

import (
	"github.com/op/go-logging"
)

type Target interface {
	Meeting(int) bool
	Exceeding(int) bool
	Delta(int) int
}

var log = logging.MustGetLogger("msstarget")
