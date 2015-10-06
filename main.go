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

	"bitbucket.org/force12io/force12-scheduler/api"
	"bitbucket.org/force12io/force12-scheduler/demand"
	"bitbucket.org/force12io/force12-scheduler/scheduler"
)

const const_sleep = 100           // milliseconds - delay before we check for demand. TODO! Make this driven by webhooks rather than simply a delay
const const_sendstate_sleep = 500 // milliseconds - delay before we send state on the metrics API

var p1TaskName string = "priority1"
var p2TaskName string = "priority2"
var p1FamilyName string = "p1-family"
var p2FamilyName string = "p2-family"

var tasks map[string]demand.Task

// cleanup resets demand for all tasks to 0 before we quit
func cleanup(s scheduler.Scheduler, tasks map[string]demand.Task) {
	var err error

	log.Println("Cleaning up tasks on interrupt")
	for name, task := range tasks {
		task.Demand = 0
		// Don't change our own force12 task
		if name != "force12" {
			err = s.StopStartNTasks(name, &task)
			if err != nil {
				log.Printf("Failed to cleanup %s tasks. %v", name, err)
				break
			}
			log.Printf("Reset %s tasks to 0 for cleanup", name)
		}
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

	// Loop, continually checking for changes in demand that need to be scheduled
	// At the moment we plough on regardless in the face of errors, simply logging them out
	for {
		select {
		case <-timeout.C:
			// Don't change demand more often than defined by demandInterval
			// We check for changes in demand more often because we want to react quickly if there hasn't been a recent change
			if time.Since(lastDemandUpdate) > st.demandInterval {
				err = handleDemandChange(di, s, tasks)
				if err != nil {
					log.Printf("Failed to handle demand change. %v", err)
				}
				lastDemandUpdate = time.Now()
			}

		case <-sendstateTimeout.C:
			// Find out how many isntances of each task are running
			s.CountAllTasks(tasks)
			err = api.SendState(st.userID, tasks, st.maxContainers)
			if err != nil {
				log.Printf("Failed to send state. %v", err)
			}

		case <-closedown:
			cleanup(s, tasks)
			os.Exit(1)
		}
	}
}
