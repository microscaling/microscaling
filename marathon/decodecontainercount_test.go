package marathon

import (
	"testing"
)

func TestDecodeContainerCount(t *testing.T) {
	count, err := DecodeContainerCount("OQ==")
	if err != nil {
		t.Fatalf("Error returned. %v", err)
	}
	if count != 9 {
		t.Fatalf("Decoded count not as expected. Have %d", count)
	}
}
