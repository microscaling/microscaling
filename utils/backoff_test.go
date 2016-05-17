package utils

import (
	"testing"
)

func TestBackoff(t *testing.T) {
	b := &Backoff{
		Factor: 10,
	}

	c := make(chan struct{}, 1)
	b.Backoff(c)

	if !b.Waiting() {
		t.Fatal("Backoff unexpectedly not waiting")
	}

	b.Stop()
	if b.Waiting() {
		t.Fatal("Backoff unexpectedly waiting")
	}

	b.Backoff(c)
	<-c

	if b.Waiting() {
		t.Fatal("Backoff unexpectedly waiting")
	}

}
