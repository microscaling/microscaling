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
	"github.com/microscaling/microscaling/utils"
)

const constGetMetricsTimeout = 500  // milliseconds - read state from the scheduler this often
const constSendMetricsTimeout = 500 // milliseconds - send on the metrics API this often

var (
	log = logging.MustGetLogger("mssagent")
)

func init() {
	initLogging()
}

// cleanup resets demand for all tasks to 0 before we quit
func cleanup(s scheduler.Scheduler, tasks *demand.Tasks) {
	tasks.Lock()
	for _, task := range tasks.Tasks {
		task.Demand = 0
	}
	tasks.Unlock()

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

	// Sending an empty struct on this channel triggers the scheduler to make updates
	demandUpdate := make(chan struct{}, 1)

	s, err := getScheduler(st, demandUpdate)
	if err != nil {
		log.Errorf("Failed to get scheduler: %v", err)
		return
	}

	tasks, err = getTasks(st)
	if err != nil {
		log.Errorf("Failed to get tasks: %v", err)
		return
	}

	// Let the scheduler know about the task types.
	for _, task := range tasks.Tasks {
		err = s.InitScheduler(task)
		if err != nil {
			log.Errorf("Failed to start task %s: %v", task.Name, err)
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
	ws, err := utils.InitWebSocket(st.microscalingAPI)
	if err != nil {
		log.Errorf("Failed to open web socket: %v", err)
		return
	}

	de, err := getDemandEngine(st, ws)
	if err != nil {
		log.Errorf("Failed to get demand engine: %v", err)
		return
	}

	go de.GetDemand(tasks, demandUpdate)

	// Handle demand updates
	go func() {
		for range demandUpdate {
			err = s.StopStartTasks(tasks)
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
		}
	}()

	// Periodically send metrics to the server
	sendMetricsTimeout := time.NewTicker(constSendMetricsTimeout * time.Millisecond)
	if st.sendMetrics {
		go func() {
			for _ = range sendMetricsTimeout.C {
				err = api.SendMetrics(ws, st.userID, tasks)
				if err != nil {
					log.Errorf("Failed to send metrics. %v", err)
				}
			}
		}()
	}

	// When we're asked to close down, we don't want to handle demand updates any more
	<-closedown
	log.Info("Clean up when ready")
	// Give the scheduler a chance to do any necessary cleanup
	s.Cleanup()
	// The demand engine is responsible for closing the demandUpdate channel so that we stop
	// doing scaling operations
	de.StopDemand(demandUpdate)

	exitWaitTimeout := time.NewTicker(constGetMetricsTimeout * time.Millisecond)
	for _ = range exitWaitTimeout.C {
		if tasks.Exited() {
			log.Info("All finished")
			break
		}
	}
}
