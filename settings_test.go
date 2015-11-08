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
	var testid string = "hello"
	os.Setenv("F12_USER_ID", testid)

	s = get_settings()
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
		os.Setenv("F12_SCHEDULER", test.sched)
		st := get_settings()
		_, err = get_scheduler(st)
		if err != nil && test.pass {
			t.Fatalf("Should have been able to create %s", test.sched)
		}
		if err == nil && !test.pass {
			t.Fatalf("Should not have been able to create %s", test.sched)
		}
	}
}
