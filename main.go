// Force12.io is a package that monitors demand for resource in a system and then scales and repurposes
// containers, based on agreed "quality of service" contracts, to best handle that demand within the constraints of your existing VM
// or physical infrastructure (for v1).
//
// Force12 is defined to optimize the use of existing physical and VM resources instantly. VMs cannot be scaled in real time (it takes
// several minutes) and new physical machines take even longer. However, containers can be started or stopped at sub-second speeds,
// allowing your infrastructure to adapt itself in real time to meet system demands.
//
// Force12 is aimed at effectively using the resources you have right now - your existing VMs or physical servers - by using them as
// optimally as possible.
//
// The Force12 approach is analogous to the way that a router dynamically optimises the use of a physical network. A router is limited
// by the capacity of the lines physically connected to it. Adding additional capacity is a physical process and takes time. Routers
// therefore make decisions in real time about which packets will be prioritized on a particular line based on the packet's priority
// (defined by a "quality of service" contract).
//
// For example, at times of high bandwidth usage a router might prioritize VOIP traffic over web browsing in real time.
//
// Containers allow Force12 to make similar "instant" judgements on service prioritisation within your existing infrastructure. Routers
// make very simplistic judgments because they have limited time and cpu and they act at a per packet level. Force12 has the capability
// of making far more sophisticated judgements, although even fairly simple ones will still provide a significant new service.
//
// This prototype is a bare bones implementation of Force12.io that recognises only 1 demand type:
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

	"github.com/force12io/force12/api"
	"github.com/force12io/force12/demand"
	"github.com/force12io/force12/scheduler"
)

const const_sleep = 100           // milliseconds - delay before we check for demand. TODO! Make this driven by webhooks rather than simply a delay
const const_sendstate_sleep = 500 // milliseconds - delay before we send state on the metrics API

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

// For this simple prototype, Force12.io sits in a loop checking for demand changes every X milliseconds
func main() {
	var err error

	st := get_settings()
	di, err := get_demand_input(st)
	if err != nil {
		log.Printf("Failed to get demand input: %v", err)
		return
	}

	s, err := get_scheduler(st)
	if err != nil {
		log.Printf("Failed to get scheduler: %v", err)
		return
	}

	tasks := get_tasks(st)
	log.Printf("Vary tasks with delta %d up to max %d containers every %d s", st.demandDelta, st.maxContainers, int(st.demandInterval.Seconds()))

	// Let the scheduler know about the task types. For the moment the actual container information is hard-coded
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

	// Periodically check for changes in demand
	// TODO: In future demand changes should come in through a channel rather than on a timer
	timeout := time.NewTicker(const_sleep * time.Millisecond)
	lastDemandUpdate := time.Now()

	// Periodically send state to the API if required
	var sendstateTimeout *time.Ticker
	if st.sendstate {
		sendstateTimeout = time.NewTicker(const_sendstate_sleep * time.Millisecond)
	}

	// Only allow one scaling command and one Send State API call to be outstanding at a time
	ready := make(chan struct{}, 1)
	ss_ready := make(chan struct{}, 1)
	var scaling_ready bool = true
	var sendState_ready bool = true
	var cleanup_when_ready bool = false
	var exit_when_ready bool = false

	// Loop, continually checking for changes in demand that need to be scheduled
	// At the moment we plough on regardless in the face of errors, simply logging them out
	for {
		select {
		case <-timeout.C:
			// Don't do anything if we're about to exit
			if cleanup_when_ready || exit_when_ready {
				break
			}

			// Don't change demand more often than defined by demandInterval
			// We check for changes in demand more often because we want to react quickly if there hasn't been a recent change
			if time.Since(lastDemandUpdate) > st.demandInterval {
				// If we already have a scaling change outstanding, we can't do another one
				if !scaling_ready {
					log.Printf("Scale change still outstanding - demand changes coming too fast to handle!")
					// This isn't an error - we simply don't try to update scale until the scheduler is ready
				} else {
					scaling_ready = false
					go func() {
						err = handleDemandChange(di, s, tasks)
						if err != nil {
							log.Printf("Failed to handle demand change. %v", err)
						}
						lastDemandUpdate = time.Now()

						// Notify the channel when the scaling command has completed
						ready <- struct{}{}
					}()
				}
			}

		case <-sendstateTimeout.C:
			// Find out how many instances of each task are running
			err = s.CountAllTasks(tasks)
			if err != nil {
				log.Printf("Failed to count containers. %v", err)
			}

			if !sendState_ready {
				log.Printf("Send state change still outstanding - can't send again yet!")
				// This isn't an error - we simply don't try to send another API call until the last response comes back
			} else {
				sendState_ready = false
				go func() {
					err = api.SendState(st.userID, tasks, st.maxContainers)
					if err != nil {
						log.Printf("Failed to send state. %v", err)
					}

					// Notify the channel when the API call has completed
					ss_ready <- struct{}{}
				}()
			}

		case <-ready:
			if exit_when_ready {
				log.Printf("All finished")
				os.Exit(1)
			}
			// An outstanding scale command has finished so we are OK to send another one
			if cleanup_when_ready {
				log.Printf("Scale command finished - now we can start cleaning up")
				exit_when_ready = true
				go func() {
					cleanup(s, tasks)
					ready <- struct{}{}
				}()
			} else {
				scaling_ready = true
			}

		case <-ss_ready:
			// An outstanding API call sending state has finished so we are OK to send another one
			sendState_ready = true

		case <-closedown:
			cleanup_when_ready = true

			if scaling_ready {
				// Trigger it now
				log.Printf("Closing down - start cleanup")
				ready <- struct{}{}
			} else {
				log.Printf("Closing down - wait till we've completed the previous scale command")
			}
		}
	}
}
