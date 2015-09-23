package compose

import (
	"testing"
)

func TestComposeScheduler(t *testing.T) {
	c := NewScheduler()

	c.InitScheduler("anything")
}
