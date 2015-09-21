package marathon

import (
	"testing"
)

func TestDecodeContainerCount(t *testing.T) {
	m := NewMarathonScheduler()

	count, err := m.DecodeContainerCount("OQ==")
	if err != nil {
		t.Fatalf("Error returned. %v", err)
	}
	if count != 9 {
		t.Fatalf("Decoded count not as expected. Have %d", count)
	}
}
