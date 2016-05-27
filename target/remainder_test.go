package target

import (
	"testing"
)

func TestRemainder(t *testing.T) {
	r := NewRemainderTarget(100)

	if r.maxContainers != 100 {
		t.Fatalf("Bad max containers")
	}

	if r.Meeting(10) {
		t.Fatalf("Remainder should never meet target")
	}

	if r.Exceeding(10) {
		t.Fatalf("Remainder should never exceed target")
	}

	if r.Delta(10) != 100 {
		t.Fatalf("Remainder delta should always be maxContainers")
	}
}
