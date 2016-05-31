package config

import (
	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/metric"
	"github.com/microscaling/microscaling/target"
)

// HardcodedConfig is used for testing
type HardcodedConfig struct{}

// compile-time assert that we implement the right interface
var _ Config = (*HardcodedConfig)(nil)

// NewHardcodedConfig gets a new hardcoded config
func NewHardcodedConfig() *HardcodedConfig {
	return &HardcodedConfig{}
}

// GetApps returns hardcoded task config
func (c *HardcodedConfig) GetApps(userID string) (tasks []*demand.Task, maxContainers int, err error) {

	task := &demand.Task{
		Name:            "consumer",
		Image:           "microscaling/queue-demo:latest",
		Priority:        1,
		MinContainers:   1,
		MaxContainers:   10,
		MaxDelta:        9,
		IsScalable:      true,
		PublishAllPorts: true,
		NetworkMode:     "host",
		Target:          target.NewQueueLengthTarget(50),
		Metric:          metric.NewNSQMetric("microscaling-demo", "microscaling-demo"),
	}

	tasks = append(tasks, task)

	task = &demand.Task{
		Name:            "remainder",
		Image:           "microscaling/priority-2:latest",
		Priority:        2,
		MinContainers:   1,
		MaxContainers:   10,
		MaxDelta:        9,
		IsScalable:      true,
		PublishAllPorts: true,
		NetworkMode:     "host",
		Target:          target.NewRemainderTarget(10),
		Metric:          metric.NewNullMetric(),
	}

	tasks = append(tasks, task)

	maxContainers = 10

	return
}
