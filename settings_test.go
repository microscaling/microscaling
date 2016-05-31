package main

import (
	// "log"
	"os"
	// "strconv"
	"testing"
	// "time"
)

func TestSettings(t *testing.T) {
	var s settings
	var testid = "hello"
	os.Setenv("MSS_USER_ID", testid)

	s = getSettings()
	if s.userID != testid {
		t.Fatalf("Unexpected userID")
	}

}

func TestInitScheduler(t *testing.T) {
	var err error

	tests := []struct {
		sched string
		pass  bool
	}{
		{sched: "COMPOSE", pass: false},
		{sched: "DOCKER", pass: true},
		{sched: "ECS", pass: false},
		{sched: "KUBERNETES", pass: false},
		{sched: "MESOS", pass: false},
		{sched: "NOMAD", pass: false},
		{sched: "TOY", pass: true},
		{sched: "BLAH", pass: false},
	}

	for _, test := range tests {
		os.Setenv("MSS_SCHEDULER", test.sched)
		st := getSettings()
		_, err = getScheduler(st, nil)
		if err != nil && test.pass {
			t.Fatalf("Should have been able to create %s", test.sched)
		}
		if err == nil && !test.pass {
			t.Fatalf("Should not have been able to create %s", test.sched)
		}
	}
}

func TestInitConfig(t *testing.T) {
	tests := []struct {
		config string
		pass   bool
	}{
		{config: "FILE", pass: false},
		// {config: "SERVER", pass: true}, Need to mock out the server for this test
		{config: "HARDCODED", pass: true},
		{config: "Blah", pass: false},
	}

	for _, test := range tests {
		os.Setenv("MSS_CONFIG", test.config)
		st := getSettings()
		tasks, err := getTasks(st)
		if err != nil && test.pass {
			t.Fatalf("Should have been able to create %s", test.config)
		}
		if err == nil && !test.pass {
			t.Fatalf("Should not have been able to create %s", test.config)
		}
		if test.config == "HARDCODED" {
			if len(tasks.Tasks) != 2 {
				t.Fatal("Expected two hardcoded tasks")
			}
		}
	}

}
