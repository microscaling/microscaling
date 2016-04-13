// Microscaling is a package that monitors demand for resource in a system and then scales and repurposes
// containers, based on agreed "quality of service" contracts, to best handle that demand within the constraints of your existing VM
// or physical infrastructure (for v1).
//
// Microscaling is defined to optimize the use of existing physical and VM resources instantly. VMs cannot be scaled in real time (it takes
// several minutes) and new physical machines take even longer. However, containers can be started or stopped at sub-second speeds,
// allowing your infrastructure to adapt itself in real time to meet system demands.
//
// Microscaling is aimed at effectively using the resources you have right now - your existing VMs or physical servers - by using them as
// optimally as possible.
//
// The microscaling approach is analogous to the way that a router dynamically optimises the use of a physical network. A router is limited
// by the capacity of the lines physically connected to it. Adding additional capacity is a physical process and takes time. Routers
// therefore make decisions in real time about which packets will be prioritized on a particular line based on the packet's priority
// (defined by a "quality of service" contract).
//
// For example, at times of high bandwidth usage a router might prioritize VOIP traffic over web browsing in real time.
//
// Containers allow microscaling to make similar "instant" judgements on service prioritisation within your existing infrastructure. Routers
// make very simplistic judgments because they have limited time and cpu and they act at a per packet level. Microscaling has the capability
// of making far more sophisticated judgements, although even fairly simple ones will still provide a significant new service.
//
// This prototype is a bare bones implementation of microscaling that recognises only 1 demand type:
// randomised demand for a priority 1 service. Resources are allocated to meet this demand for priority 1, and spare resource can
// be used for a priority 2 service.
//
// These demand type examples have been chosen purely for simplicity of demonstration. In the future more demand types
// will be offered
package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"

	"github.com/microscaling/microscaling/api"
	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/scheduler"
)

const constGetMetricsTimeout = 500 // milliseconds - delay before we send read state (and optionally send on the metrics API)

var (
	log = logging.MustGetLogger("mssagent")
)

func init() {
	initLogging()
}

// cleanup resets demand for all tasks to 0 before we quit
func cleanup(s scheduler.Scheduler, running *demand.Tasks) {
	running.Lock()

	tasks := running.Tasks
	for name, task := range tasks {
		task.Demand = 0
		tasks[name] = task
	}

	running.Unlock()

	log.Debugf("Reset tasks to 0 for cleanup")
	err := s.StopStartTasks(tasks)
	if err != nil {
		log.Errorf("Failed to cleanup tasks. %v", err)
	}
}

// For this simple prototype, Microscaling sits in a loop checking for demand changes every X milliseconds
func main() {
	var err error
	var tasks *demand.Tasks

	st := getSettings()

	s, err := getScheduler(st)
	if err != nil {
		log.Errorf("Failed to get scheduler: %v", err)
		return
	}

	tasks = getTasks(st)

	// Let the scheduler know about the task types.
	for name, task := range tasks.Tasks {
		err = s.InitScheduler(name, &task)
		if err != nil {
			log.Errorf("Failed to start task %s: %v", name, err)
			return
		}
	}

	// Check if there are already any of these containers running
	err = s.CountAllTasks(tasks)
	if err != nil {
		log.Errorf("Failed to count containers. %v", err)
	}

	// Set the initial requested counts to match what's running
	for name, task := range tasks.Tasks {
		task.Requested = task.Running
		tasks.Tasks[name] = task
	}

	// Prepare for cleanup when we receive an interrupt
	closedown := make(chan os.Signal, 1)
	signal.Notify(closedown, os.Interrupt)
	signal.Notify(closedown, syscall.SIGTERM)

	// Open a web socket to the server TODO!! This won't always be necessary if we're not sending metrics & calculating demand locally
	ws, err := api.InitWebSocket()
	if err != nil {
		log.Errorf("Failed to open web socket: %v", err)
		return
	}

	demandUpdate := make(chan struct{}, 1)
	de, err := getDemandEngine(st, ws)
	if err != nil {
		log.Errorf("Failed to get demand engine: %v", err)
		return
	}

	go de.GetDemand(tasks, demandUpdate)

	// Handle demand updates
	go func() {
		for range demandUpdate {
			log.Debug("Demand update")
			tasks.Lock()
			err = s.StopStartTasks(tasks.Tasks)
			tasks.Unlock()
			if err != nil {
				log.Errorf("Failed to stop / start tasks. %v", err)
			}
		}

		// When the demandUpdate channel is closed, it's time to scale everything down to 0
		cleanup(s, tasks)
	}()

	// Periodically read the current state of tasks
	getMetricsTimeout := time.NewTicker(constGetMetricsTimeout * time.Millisecond)
	go func() {
		for _ = range getMetricsTimeout.C {
			// Find out how many instances of each task are running
			err = s.CountAllTasks(tasks)
			if err != nil {
				log.Errorf("Failed to count containers. %v", err)
			}

			if st.sendMetrics {
				log.Debug("Sending metrics")
				err = api.SendMetrics(ws, st.userID, tasks.Tasks)
				if err != nil {
					log.Errorf("Failed to send metrics. %v", err)
				}
			}
		}
	}()

	// When we're asked to close down, we don't want to handle demand updates any more
	<-closedown
	log.Info("Clean up when ready")
	// The demand engine is responsible for closing the demandUpdate channel so that we stop
	// doing scaling operations
	de.StopDemand(demandUpdate)

	exitWaitTimeout := time.NewTicker(constGetMetricsTimeout * time.Millisecond)
	for _ = range exitWaitTimeout.C {
		if demand.Exited(tasks) {
			log.Info("All finished")
			break
		}
	}
}
