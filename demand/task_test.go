package demand

import (
	"testing"

	"github.com/microscaling/microscaling/metric"
	"github.com/microscaling/microscaling/target"
)

func TestIsRemainder(t *testing.T) {
	r := target.NewRemainderTarget(10)
	q := target.NewQueueLengthTarget(10)
	sq := target.NewSimpleQueueLengthTarget(10)

	testTask := Task{
		Target: r,
	}

	if !testTask.IsRemainder() {
		t.Fatalf("Didn't recognise a remainder task")
	}

	testTask.Target = q
	if testTask.IsRemainder() {
		t.Fatalf("Mistook a queue task for a remainder")
	}

	testTask.Target = sq
	if testTask.IsRemainder() {
		t.Fatalf("Mistook a simple queue task for a remainder")
	}
}

func getTestTask() (tm *metric.ToyMetric, testTask Task) {
	tm = metric.NewToyMetric()
	m := tm
	q := target.NewSimpleQueueLengthTarget(10)

	testTask = Task{
		Target:        q,
		Metric:        m,
		MinContainers: 1,
		MaxContainers: 5,
		MaxDelta:      2,
		IsScalable:    true,
	}

	return tm, testTask
}

func TestScaleUpCount(t *testing.T) {
	m, testTask := getTestTask()

	// If the requested and ideal number of containers match, and we're meeting the target we don't scale at all
	testTask.IdealContainers = 2
	testTask.Requested = 2
	if testTask.ScaleUpCount() != 0 {
		t.Fatalf("Unexpected scale up count (0)")
	}

	// Now we're not meeting the target so we expect some scaling
	m.SettableCurrent = 100

	// Never scale to more than max containers
	testTask.IdealContainers = 10
	testTask.Requested = 4
	if testTask.ScaleUpCount() != 1 {
		t.Fatalf("Unexpected scale up count (1)")
	}

	// Never scale by more than max delta
	testTask.IdealContainers = 10
	testTask.Requested = 2
	if testTask.ScaleUpCount() != 2 {
		t.Fatalf("Unexpected scale up count (2)")
	}

	// Scale up count should be non-negative
	testTask.IdealContainers = 1
	if testTask.ScaleUpCount() != 0 {
		t.Fatalf("Unexpected scale up count (3)")
	}

	// Normal scale up case
	testTask.IdealContainers = 4
	if testTask.ScaleUpCount() != 2 {
		t.Fatalf("Unexpected scale up count (4)")
	}
}

func TestScaleDownCount(t *testing.T) {
	m, testTask := getTestTask()

	// If the requested and ideal number of containers match, and we're meeting the target we don't scale at all
	testTask.IdealContainers = 2
	testTask.Requested = 2
	if testTask.ScaleDownCount() != 0 {
		t.Fatalf("Unexpected scale down count (0)")
	}

	// Now we're not meeting the target so we expect some scaling
	m.SettableCurrent = 1

	// Never scale to fewer than min containers
	testTask.IdealContainers = 0
	testTask.Requested = 2
	if testTask.ScaleDownCount() != -1 {
		t.Fatalf("Unexpected scale down count (1)")
	}

	// Never scale by more than max delta
	testTask.IdealContainers = 1
	testTask.Requested = 5
	if testTask.ScaleDownCount() != -2 {
		t.Fatalf("Unexpected scale down count (2)")
	}

	// Scale down count should be negative
	testTask.IdealContainers = 10
	if testTask.ScaleDownCount() != 0 {
		t.Fatalf("Unexpected scale down count (3)")
	}

	// Normal scale down case
	testTask.IdealContainers = 4
	if testTask.ScaleDownCount() != -1 {
		t.Fatalf("Unexpected scale down count (4)")
	}
}

func TestCanScaleDown(t *testing.T) {
	testTask := &Task{
		IsScalable:    false,
		MinContainers: 3,
		Requested:     3,
	}

	if testTask.CanScaleDown() != 0 {
		t.Fatalf("Can't scale down an unscalable task")
	}

	testTask.IsScalable = true
	if testTask.CanScaleDown() != 0 {
		t.Fatalf("Can't scale down if requested is already at minimum")
	}

	testTask.Requested = 5
	if testTask.CanScaleDown() != 2 {
		t.Fatalf("Can't scale down if requested is already at minimum")
	}
}
