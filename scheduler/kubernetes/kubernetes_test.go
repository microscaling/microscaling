package kubernetes

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/microscaling/microscaling/demand"
)

func TestKubernetesInitScheduler(t *testing.T) {
	kubeConfig, _ := filepath.Abs("kube_config_test")
	namespace := "default"
	demandUpdate := make(chan struct{}, 1)

	k := NewScheduler(kubeConfig, namespace, demandUpdate)

	task := demand.Task{
		Name:   "consumer",
		Demand: 5,
	}

	k.InitScheduler(&task)

	if k.namespace != namespace {
		t.Errorf("Expected namespace to be %s but was %s", namespace, k.namespace)
	}

	if k.backoff.Factor != 2 {
		t.Errorf("Expected backoff factor to be 2 but was %d", k.backoff.Factor)
	}

	if k.backoff.Min != (250 * time.Millisecond) {
		t.Errorf("Expected min backoff to be 250 ms but was %d", k.backoff.Min)
	}

	if k.backoff.Max != (5 * time.Second) {
		t.Errorf("Expected max backoff to be 5 secs but was %d", k.backoff.Max)
	}
}
