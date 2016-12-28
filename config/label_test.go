package config

import (
	"testing"

	"github.com/microscaling/microscaling/demand"
)

func TestLabelConfig(t *testing.T) {
	var task demand.Task
	labels := make(map[string]string, 5)

	labels["com.microscaling.is-scalable"] = "True"
	labels["com.microscaling.PRIoRiTY"] = "5" // Check we can cope if it's not lower case
	labels["com.microscaling.max-delta"] = "2"
	labels["com.microscaling.min-containers"] = "1"
	labels["com.microscaling.MAX-containers"] = "20"

	parseLabels(&task, labels)

	if !task.IsScalable {
		t.Errorf("Bad IsScalable")
	}

	if task.Priority != 5 {
		t.Errorf("Bad Priority")
	}

	if task.MaxDelta != 2 {
		t.Errorf("Bad Max Delta")
	}

	if task.MinContainers != 1 {
		t.Errorf("Bad Min Containers")
	}

	if task.MaxContainers != 20 {
		t.Errorf("Bad Min Containers")
	}

}
