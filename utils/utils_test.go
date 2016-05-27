package utils

import (
	"os"
	"testing"
)

func TestEnvFl64(t *testing.T) {

	e := EnvFl64("MSS_A_VAR_WE_WOULDNT_SET", 0.7)
	if e != 0.7 {
		t.Fatalf("unexpected default value %f", e)
	}

	os.Setenv("MSS_A_VAR_WE_WOULDNT_SET", "0.5")
	e = EnvFl64("MSS_A_VAR_WE_WOULDNT_SET", 0.7)
	if e != 0.5 {
		t.Fatalf("unexpected set value %f", e)
	}

	os.Setenv("MSS_A_VAR_WE_WOULDNT_SET", "XX")
	e = EnvFl64("MSS_A_VAR_WE_WOULDNT_SET", 0.7)
	if e != 0.7 {
		t.Fatalf("unexpected default overriding bad value %f", e)
	}
}
