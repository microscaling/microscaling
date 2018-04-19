package config

import (
	"testing"

	"github.com/microscaling/microscaling/demand"
)

func TestEnvVarConfig(t *testing.T) {
	tests := []struct {
		name          string
		json          string
		success       bool
		tasks         []*demand.Task
		maxContainers int
	}{
		{
			name: "basic match",
			json: `
			{
				"name": "world",
				"maxContainers": 10,
				"apps": [
					{
						"name": "priority1",
						"priority": 1
					},
					{
						"name": "priority2",
						"priority": 2
					}
				]
			}`,
			success: true,
			tasks: []*demand.Task{
				&demand.Task{
					Name:     "priority1",
					Priority: 1,
				},
				&demand.Task{
					Name:     "priority2",
					Priority: 2,
				},
			},
			maxContainers: 10,
		},
		{
			name:    "empty json",
			json:    "",
			success: false,
		},
		{
			name: "invalid json",
			json: `
			{
				"name": "invalid",
				"maxContainers": "Not a number"
			}`,
			success: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := NewEnvVarConfig(tc.json)
			tasks, maxC, err := c.GetApps("test-user")

			if tc.success {
				if err != nil {
					t.Fatalf("Expected error to be nil but was %v", err)
				}
				if len(tc.tasks) != len(tasks) {
					t.Fatalf("Expected %d tasks but was %d", len(tc.tasks), len(tasks))
				}

				for i, task := range tasks {
					et := tc.tasks[i]

					if et.Name != task.Name {
						t.Fatalf("Expected task %d name to be %s but was %s", i, et.Name, task.Name)
					}
					if et.Priority != task.Priority {
						t.Fatalf("Expected task %d priority to be %d but was %d", i, et.Priority, task.Priority)
					}
				}

				if tc.maxContainers != maxC {
					t.Fatalf("Expected max containers to be %d but was %d", tc.maxContainers, maxC)
				}

			} else {
				if err == nil {
					t.Fatalf("Expected an error but was %v", err)
				}
			}
		})
	}
}
