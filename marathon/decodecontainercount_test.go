package marathon

import (
	"testing"
)

func TestDecodeContainerCount(t *testing.T) {

	text := `[
   {
       "CreateIndex": 8,
       "ModifyIndex": 15,
       "LockIndex": 0,
       "Key": "priority1-demand",
       "Flags": 0,
       "Value": "OQ=="
   }
]`
	count := DecodeContainerCount(text)
	if count != 9 {
		t.Fatalf("Decoded count not as expected. Have %d", count)
	}
}
