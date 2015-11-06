package demandapi

import (
	"log"

	"github.com/force12io/force12/api"
	"github.com/force12io/force12/demand"
)

type DemandFromApi struct {
	userId string
}

// check that we implement the demand interface
var _ demand.Input = (*DemandFromApi)(nil)

func NewDemandModel(userId string) *DemandFromApi {

	return &DemandFromApi{
		userId: userId,
	}
}

func (i *DemandFromApi) Update(tasks map[string]demand.Task) (demandChanged bool, err error) {
	demandChanged = false

	td, err := api.GetTasks(i.userId)
	if err != nil {
		log.Printf("Problem getting tasks: %v", err)
		return
	}

	for _, task := range td {
		name := task.App

		if existing_task, ok := tasks[name]; ok {
			if existing_task.Demand != task.DemandCount {
				demandChanged = true
			}
			existing_task.Demand = task.DemandCount
			tasks[name] = existing_task
		}
	}

	return
}
