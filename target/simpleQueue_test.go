package target

import (
	"testing"
)

func TestNewSimpleQueueTarget(t *testing.T) {
	length := 10
	q := NewSimpleQueueLengthTarget(length)

	if q.length != length {
		t.Fatalf("Wrong length")
	}

	minLength := float64(length) * queueLengthExceedingPercent
	if q.minLength != int(minLength) {
		t.Fatalf("Wrong minLength")
	}

}

func TestSimpleQueue(t *testing.T) {
	target := 10
	sq := NewSimpleQueueLengthTarget(target)
	queueTargetTest(t, sq)

	d := sq.Delta(20)
	if d != 1 {
		t.Fatalf("Bad delta")
	}

	d = sq.Delta(5)
	if d != -1 {
		t.Fatalf("Bad -ve delta")
	}
}
