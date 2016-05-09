package api

// DemandPayload is the JSON sent from the server describing the number of containers needed for each task
// This is only used if we're generating demand server-sde
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
