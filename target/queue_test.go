package target

import (
	"testing"
)

func TestNewQueueTarget(t *testing.T) {
	length := 10
	q := NewQueueLengthTarget(length)

	if q.length != length {
		t.Fatalf("Wrong length")
	}

	minLength := float64(length) * queueLengthExceedingPercent
	if q.minLength != int(minLength) {
		t.Fatalf("Wrong minLength")
	}

	if q.velSamples != 1 {
		t.Fatalf("Wrong velocity samples")
	}
}

func queueTargetTest(t *testing.T, q Target) {
	meeting := 10
	notMeeting := 11
	exceeding := 5
	notExceeding := 9

	if !q.Meeting(meeting) {
		t.Fatalf("Not meeting")
	}

	if !q.Exceeding(exceeding) {
		t.Fatalf("Not exceeding")
	}

	// Exceeding should always be meeting
	if !q.Meeting(exceeding) {
		t.Fatalf("Not meeting with exceeding")
	}

	if q.Meeting(notMeeting) {
		t.Fatalf("Unexpectedly meeting")
	}

	if q.Exceeding(notExceeding) {
		t.Fatalf("Unexpectedly exceeding")
	}

	// In this example we expect notExceeding to me meeting
	if !q.Meeting(notExceeding) {
		t.Fatalf("Unexpectedly not meeting with notExceeding")
	}
}

func TestQueue(t *testing.T) {
	target := 10
	q := NewQueueLengthTarget(target)
	queueTargetTest(t, q)

	q.kP = 1
	q.kD = 1
	q.kI = 0

	// First time we can't measure velocity
	d := q.Delta(20)
	if d != 10 {
		t.Fatalf("Bad delta (0)")
	}

	d = q.Delta(30)
	// Err = 20; vel = 10 -> kD * err + kP * vel = 30
	if d != 30 {
		t.Fatalf("Bad delta (1)")
	}

	d = q.Delta(30)
	// Err = 20; vel = 0 -> kD * err + kP * vel = 20
	if d != 20 {
		t.Fatalf("Bad delta (2)")
	}

	d = q.Delta(5)
	// Err = -5; vel = -25 -> kD * err + kP * vel = -30
	if d != -30 {
		t.Fatalf("Bad delta (3)")
	}
}
