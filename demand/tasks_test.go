package demand

import (
	"testing"
)

func getTestTasks() Tasks {
	tt := Tasks{
		MaxContainers: 10,
		Tasks:         make([]*Task, 3),
	}

	tt.Tasks[0] = &Task{
		Name:      "Zero",
		Requested: 2,
		Priority:  0,
		Running:   2,
	}

	tt.Tasks[1] = &Task{
		Name:      "One",
		Requested: 2,
		Priority:  1,
		Running:   2,
	}

	tt.Tasks[2] = &Task{
		Name:      "Two",
		Requested: 2,
		Priority:  2,
		Running:   2,
	}

	return tt
}

func TestCheckCapacity(t *testing.T) {

	tt := getTestTasks()

	// Max of 10, currently 6 requested
	if tt.CheckCapacity() != 4 {
		t.Fatalf("Bad capcity check")
	}
}

func TestGetName(t *testing.T) {
	tt := getTestTasks()

	_, err := tt.GetTask("Zero")
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = tt.GetTask("missing")
	if err == nil {
		t.Fatalf("Should have failed to get non-existant task")
	}
}

func TestPrioritySort(t *testing.T) {
	tt := getTestTasks()

	tt.PrioritySort(false)
	for _, v := range tt.Tasks {
		if v.Priority != 0 {
			t.Fatalf("Badly sorted in normal order: first is %v", v)
		}
		break
	}

	tt.PrioritySort(true)
	for _, v := range tt.Tasks {
		if v.Priority != 2 {
			t.Fatalf("Badly sorted in reverse order: first is %v", v)
		}
		break
	}
}

func TestExited(t *testing.T) {
	tt := getTestTasks()

	if tt.Exited() {
		t.Fatal("Unexpectedly we look like we exited")
	}

	for _, v := range tt.Tasks {
		v.Running = 0
	}

	if !tt.Exited() {
		t.Fatal("Unexpectedly not exited")
	}
}
