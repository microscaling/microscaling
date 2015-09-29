// Force12.io is a package that monitors demand for resource in a system and then scales and repurposes
// containers, based on agreed "quality of service" contracts, to best handle that demand within the constraints of your existing VM
// or physical infrastructure (for v1).
//
// Force12 is defined to optimize the use of existing physical and VM resources instantly. VMs cannot be scaled in real time (it takes
// several minutes) and new physical machines take even longer. However, containers can be started or stopped at sub second speeds,
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
//
// V1 - Force12.io reacts to increased demand by starting/stopping containers on the slaves already in play.
//
// This version of Force12 starts and stops containers on a Mesos cluser using Marathon as the scheduler
package main

import (
	"log"
	"os"
	"time"

	"bitbucket.org/force12io/force12-scheduler/api"
	"bitbucket.org/force12io/force12-scheduler/compose"
	"bitbucket.org/force12io/force12-scheduler/consul"
	"bitbucket.org/force12io/force12-scheduler/demand"
	"bitbucket.org/force12io/force12-scheduler/marathon"
	"bitbucket.org/force12io/force12-scheduler/rng"
	"bitbucket.org/force12io/force12-scheduler/scheduler"
	"bitbucket.org/force12io/force12-scheduler/toy_scheduler"
)

const const_sleep = 100           // milliseconds
const const_sendstate_sleeps = 5  // number of sleeps before we send state on the API
const const_stopsleep = 250       // milliseconds pause between stopping and restarting containers
const const_p1demandstart int = 1 // The yaml file will automatically start one of each
const const_p2demandstart int = 1
const const_maxcontainers int = 9

var p1TaskName string = "priority1"
var p2TaskName string = "priority2"
var p1FamilyName string = "p1-family"
var p2FamilyName string = "p2-family"

var tasks map[string]demand.Task

// handleDemandChange checks the new demand
func handleDemandChange(input demand.Input, s scheduler.Scheduler) error {
	var err error = nil
	var demandChanged bool
	var t demand.Task

	demandChanged, err = update(input)
	if err != nil {
		log.Printf("Failed to get new demand. %v", err)
		return err
	}

	if demandChanged {
		// Ask the scheduler to make the changes
		for name, task := range tasks {
			t = task
			err = s.StopStartNTasks(name, task.FamilyName, t.Demand, &t.Requested)
			if err != nil {
				log.Printf("Failed to start %s tasks. %v", name, err)
				break
			}

			tasks[name] = t
		}
	}

	return err
}

// update checks for changes in demand, returning true if demand changed
func update(input demand.Input) (bool, error) {
	var err error = nil
	var demandchange bool

	// TODO! Make this less tied to the p1 / p2 simple model
	var p1 demand.Task = tasks[p1TaskName]
	var p2 demand.Task = tasks[p2TaskName]

	// Save the old demand
	oldP1Demand := p1.Demand
	oldP2Demand := p2.Demand

	p1.Demand, err = input.GetDemand(p1TaskName)
	if err != nil {
		log.Printf("Failed to get new demand for task %s. %v", p1TaskName, err)
		return false, err
	}
	//log.Printf("container count %v\n", container_count)
	p2.Demand = const_maxcontainers - p1.Demand

	//Has the demand changed?
	demandchange = (p1.Demand != oldP1Demand) || (p2.Demand != oldP2Demand)

	// Update tsaks map
	tasks[p1TaskName] = p1
	tasks[p2TaskName] = p2

	// This is where we could decide whether this is a significant enough
	// demand change to do anything

	log.Println(tasks)

	return demandchange, err
}

func getEnvOrDefault(name string, defaultValue string) string {
	v := os.Getenv(name)
	if v == "" {
		v = defaultValue
	}

	return v
}

// For the simple prototype, Force12.io sits in a loop checking for demand changes every X milliseconds
// In phase 2 we'll add a reactive mode where appropriate.
//
// Note - we don't route messages from demandcheckers to demandhandlers using channels because we want new values
// to override old values. Queued history is of no importance here.
//
// Also for simplicity this first release is concurrency free (single threaded)
func main() {
	var err error
	// TODO! Make it so you can send in settings on the command line
	var demandModelType string = getEnvOrDefault("F12_DEMAND_MODEL", "RNG")
	var schedulerType string = getEnvOrDefault("F12_SCHEDULER", "COMPOSE")
	var sendstateString string = getEnvOrDefault("F12_SEND_STATE_TO_API", "true")
	var sendstate bool = (sendstateString == "true")
	p1TaskName = getEnvOrDefault("F12_PRIORITY1_TASK", p1TaskName)
	p2TaskName = getEnvOrDefault("F12_PRIORITY2_TASK", p2TaskName)
	// TODO!! FInd out what CLIENT/SERVER_FAMILY should default to
	p1FamilyName = getEnvOrDefault("F12_PRIORITY1_FAMILY", p1FamilyName)
	p2FamilyName = getEnvOrDefault("F12_PRIORITY2_FAMILY", p2FamilyName)

	var di demand.Input
	var s scheduler.Scheduler

	switch demandModelType {
	case "CONSUL":
		log.Println("Getting demand metric from Consul")
		di = consul.NewDemandModel()
	case "RNG":
		log.Println("Random demand generation")
		di = rng.NewDemandModel()
	default:
		log.Printf("Bad value for F12_DEMAND_MODEL: %s", demandModelType)
		return
	}

	switch schedulerType {
	case "COMPOSE":
		log.Println("Scheduling with Docker compose")
		s = compose.NewScheduler()
	case "ECS":
		log.Println("Scheduling with ECS not yet supported")
		return
	case "MESOS":
		log.Println("Scheduling with Mesos / Marathon")
		s = marathon.NewScheduler()
	case "TOY":
		log.Println("Scheduling with toy scheduler")
		s = toy_scheduler.NewScheduler()
	default:
		log.Printf("Bad value for F12_SCHEDULER: %s", schedulerType)
		return
	}

	// Initialise task types
	tasks = make(map[string]demand.Task)

	tasks[p1TaskName] = demand.Task{
		FamilyName: p1FamilyName,
		Demand:     const_p1demandstart,
		Requested:  0,
	}

	tasks[p2TaskName] = demand.Task{
		FamilyName: p2FamilyName,
		Demand:     const_p2demandstart,
		Requested:  0,
	}

	// Let the scheduler know about these task types. For the moment the actual container information is hard-coded
	for name, _ := range tasks {
		err = s.InitScheduler(name)
		if err != nil {
			log.Printf("Failed to start task %s: %v", name, err)
			return
		}
	}

	var sleepcount int = 0
	var sleep time.Duration = const_sleep * time.Millisecond

	// Loop, continually checking for changes in demand that need to be scheduled
	// At the moment we plough on regardless in the face of errors, simply logging them out
	for {
		err = handleDemandChange(di, s)
		if err != nil {
			log.Printf("Failed to handle demand change. %v", err)
		}

		time.Sleep(sleep)
		sleepcount++
		if sleepcount == const_sendstate_sleeps {
			sleepcount = 0

			//Periodically send state to the API if required
			if sendstate {
				// TODO! User ID needs to be not hardcoded
				err = api.SendState("5k5gk", s, tasks)
				if err != nil {
					log.Printf("Failed to send state. %v", err)
				}
			}
		}
	}
}
