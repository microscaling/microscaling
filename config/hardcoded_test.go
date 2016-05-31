package config

import (
	"testing"
)

func TestGetApps(t *testing.T) {
	c := NewHardcodedConfig()
	tasks, maxC, err := c.GetApps("hello")

	if len(tasks) != 2 {
		t.Fatal("Expected two hardcoded tasks")
	}

	if maxC != 10 {
		t.Fatal("Expected max containers of 10")
	}

	if err != nil {
		t.Fatalf("Shouldn't fail to get hardcoded config")
	}
}
