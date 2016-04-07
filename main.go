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
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/microscaling/microscaling/api"
	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/scheduler"
)

const constSendMetricsSleep = 500 // milliseconds - delay before we send state on the metrics API

var tasks map[string]demand.Task

// cleanup resets demand for all tasks to 0 before we quit
func cleanup(s scheduler.Scheduler, tasks map[string]demand.Task) {
	var err error

	for name, task := range tasks {
		task.Demand = 0
		tasks[name] = task
	}

	log.Printf("Reset tasks to 0 for cleanup")
	err = s.StopStartTasks(tasks)
	if err != nil {
		log.Printf("Failed to cleanup tasks. %v", err)
	}
}

// For this simple prototype, Microscaling sits in a loop checking for demand changes every X milliseconds
func main() {
	var err error

	st := getSettings()

	s, err := getScheduler(st)
	if err != nil {
		log.Printf("Failed to get scheduler: %v", err)
		return
	}

	tasks := getTasks(st)

	// Let the scheduler know about the task types.
	for name, task := range tasks {
		err = s.InitScheduler(name, &task)
		if err != nil {
			log.Printf("Failed to start task %s: %v", name, err)
			return
		}
	}

	// Check if there are already any of these containers running
	err = s.CountAllTasks(tasks)
	if err != nil {
		log.Printf("Failed to count containers. %v", err)
	}

	// Prepare for cleanup when we receive an interrupt
	closedown := make(chan os.Signal, 1)
	signal.Notify(closedown, os.Interrupt)
	signal.Notify(closedown, syscall.SIGTERM)

	// Listen for demand on a websocket (we'll also use this to send metrics)
	demandUpdate := make(chan []api.TaskDemand, 1)
	ws, err := api.InitWebSocket()
	go api.Listen(ws, demandUpdate)

	// Periodically send state to the API if required
	var sendMetricsTimeout *time.Ticker
	if st.sendMetrics {
		sendMetricsTimeout = time.NewTicker(constSendMetricsSleep * time.Millisecond)
	}

	// Only allow one scaling command and one metrics send to be outstanding at a time
	ready := make(chan struct{}, 1)
	metricsReady := make(chan struct{}, 1)
	var scalingReady = true
	var sendMetricsReady = true
	var cleanupWhenReady = false
	var exitWhenReady = false

	// Loop, continually checking for changes in demand that need to be scheduled
	// At the moment we plough on regardless in the face of errors, simply logging them out
	for {
		select {
		case td := <-demandUpdate:
			// Don't do anything if we're about to exit
			if cleanupWhenReady || exitWhenReady {
				break
			}

			// If we already have a scaling change outstanding, we can't do another one
			if scalingReady {
				scalingReady = false
				go func() {
					err = handleDemandChange(td, s, tasks)
					if err != nil {
						log.Printf("Failed to handle demand change. %v", err)
					}

					// Notify the channel when the scaling command has completed
					ready <- struct{}{}
				}()
			} else {
				log.Println("Scale still outstanding")
			}

		case <-sendMetricsTimeout.C:
			if sendMetricsReady {
				log.Println("Sending metrics")
				sendMetricsReady = false
				go func() {
					// Find out how many instances of each task are running
					err = s.CountAllTasks(tasks)
					if err != nil {
						log.Printf("Failed to count containers. %v", err)
					}

					err = api.SendMetrics(ws, st.userID, tasks)
					if err != nil {
						log.Printf("Failed to send metrics. %v", err)
					}

					// Notify the channel when the API call has completed
					metricsReady <- struct{}{}
				}()
			} else {
				log.Println("Not ready to send metrics")
			}

		case <-ready:
			if exitWhenReady {
				log.Printf("All finished")
				os.Exit(1)
			}

			// An outstanding scale command has finished so we are OK to send another one
			if cleanupWhenReady {
				log.Printf("Cleaning up")
				exitWhenReady = true
				go func() {
					cleanup(s, tasks)
					ready <- struct{}{}
				}()
			} else {
				scalingReady = true
			}

		case <-metricsReady:
			// Finished sending metrics so we are OK to send another one
			sendMetricsReady = true

		case <-closedown:
			log.Printf("Clean up when ready")
			cleanupWhenReady = true
			if scalingReady {
				// Trigger it now
				ready <- struct{}{}
			}
		}
	}
}
