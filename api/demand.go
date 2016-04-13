// API between Microscaling agent and server
package api

type DemandPayload struct {
	Demand DemandUpdate `json:"demand"`
}

type DemandUpdate struct {
	Tasks []TaskDemand `json:"tasks"`
}

type TaskDemand struct {
	App         string `json:"app"`
	DemandCount int    `json:"demandCount"`
}
