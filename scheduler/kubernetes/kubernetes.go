// Package kubernetes provides a scheduler using the Kubernetes API.
package kubernetes

import (
	"encoding/json"
	"time"

	"github.com/op/go-logging"

	"k8s.io/client-go/1.5/kubernetes"
	"k8s.io/client-go/1.5/pkg/api"

	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/scheduler"
	"github.com/microscaling/microscaling/utils"
)

var (
	log = logging.MustGetLogger("mssscheduler")
	pt  = api.StrategicMergePatchType
)

// KubernetesScheduler holds the Kubernetes clientset and a Backoff struct for each task.
type KubernetesScheduler struct {
	clientset    *kubernetes.Clientset
	namespace    string
	demandUpdate chan struct{}
	backoff      *utils.Backoff
}

// deploymentSpec updates the Deployments API
type deploymentSpec struct {
	Spec deployment `json:"spec"`
}

// deployment sends the number of pods to run
type deployment struct {
	Replicas int32 `json:"replicas"`
}

// NewScheduler returns a pointer to the scheduler. Creates k8s clientset from the provided kube
// config or when running as a pod uses the in cluster config.
func NewScheduler(kubeConfig string, namespace string, demandUpdate chan struct{}) *KubernetesScheduler {
	clientset, err := utils.NewKubeClientset(kubeConfig, namespace)
	if err != nil {
		log.Errorf("Error creating Kubernetes clientset: %v", err)
		return nil
	}

	return &KubernetesScheduler{
		clientset:    clientset,
		namespace:    namespace,
		demandUpdate: demandUpdate,
		backoff: &utils.Backoff{
			Min:    250 * time.Millisecond,
			Max:    5 * time.Second,
			Factor: 2,
		},
	}
}

// compile-time assert that we implement the right interface
var _ scheduler.Scheduler = (*KubernetesScheduler)(nil)

// InitScheduler initializes the scheduler.
func (k *KubernetesScheduler) InitScheduler(task *demand.Task) (err error) {
	log.Infof("Kubernetes initializing task %s", task.Name)
	return err
}

// StopStartTasks by calling the Kubernetes Deployments API.
func (k *KubernetesScheduler) StopStartTasks(tasks *demand.Tasks) error {
	// Create tasks if there aren't enough of them, and stop them if there are too many
	var tooMany []*demand.Task
	var tooFew []*demand.Task
	var err error

	// Check we're not already backed off. This could easily happen if we get a demand update
	// arrive while we are in the midst of a previous backoff.
	if k.backoff.Waiting() {
		log.Debug("Backoff timer still running")
		return nil
	}

	tasks.Lock()
	defer tasks.Unlock()

	for _, t := range tasks.Tasks {
		if t.Demand > t.Requested {
			// There aren't enough of these pods yet
			tooFew = append(tooFew, t)
		}
		if t.Demand < t.Requested {
			// There are too many of these pods
			tooMany = append(tooMany, t)
		}
	}

	// Concatentate the two lists - scale down first to free up resources
	tasksToScale := append(tooMany, tooFew...)
	for _, t := range tasksToScale {
		log.Debugf("Scaling task %s to %d", t.Name, t.Demand)

		running, err := k.countTasks(t.Name)
		if err != nil {
			log.Errorf("Error getting task count for %s: %v", t.Name, err)
			return err
		}

		if running == t.Requested {
			// Clear any backoffs before scaling
			k.backoff.Reset()

			err := k.stopStartTask(t)
			if err != nil {
				log.Errorf("Error scaling %s: %v ", t.Name, err)
				return err
			}

			log.Infof("Scaled %s to %d", t.Name, t.Demand)

		} else {
			// Trigger a backoff as the previous scaling action is not yet complete
			err = k.backoff.Backoff(k.demandUpdate)
			log.Debugf("Backing off %s %d requested but %d running", t.Name, t.Demand, running)

			return err
		}
	}

	return err
}

// CountAllTasks tells us how many pods of each deployment are currently running.
func (k *KubernetesScheduler) CountAllTasks(running *demand.Tasks) (err error) {
	running.Lock()
	defer running.Unlock()

	// Set running counts. Defaults to 0 if the deployment does not exist.
	tasks := running.Tasks
	for _, t := range tasks {
		running, err := k.countTasks(t.Name)
		if err != nil {
			log.Errorf("Error getting deployment %s: %v", t.Name, err)
			return err
		}

		t.Running = running
		log.Debugf("Deployment %s: requested %d, running %d", t.Name, t.Requested, running)
	}

	return err
}

// stopStartTask patches the deployment to set the desired number of pods
func (k *KubernetesScheduler) stopStartTask(task *demand.Task) (err error) {
	deployment := deploymentSpec{
		deployment{
			Replicas: int32(task.Demand),
		},
	}

	bytes, err := json.Marshal(deployment)
	if err != nil {
		log.Errorf("Error marshaling deployment json for %s: %v", task.Name, err)
		return err
	}

	_, err = k.clientset.Extensions().Deployments(k.namespace).Patch(task.Name, pt, bytes)
	if err != nil {
		log.Errorf("Error patching deployment %s: %v", task.Name, err)
		return err
	}

	task.Requested = task.Demand

	return err
}

// countTasks counts how many running pods exist for the deployment
func (k *KubernetesScheduler) countTasks(taskName string) (count int, err error) {
	d, err := k.clientset.Extensions().Deployments(k.namespace).Get(taskName)
	if err != nil {
		log.Errorf("Error getting deployment %s: %v", taskName, err)
		return count, err
	}

	count = int(d.Status.AvailableReplicas)
	return count, err
}

// Cleanup gives the scheduler an opportunity to stop anything that needs to be stopped
func (k *KubernetesScheduler) Cleanup() error {
	k.backoff.Stop()
	return nil
}
