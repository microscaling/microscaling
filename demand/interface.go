// The demand package defined the interface for demand models
package demand

type Task struct {
	Demand     int
	Requested  int
	FamilyName string
}

type Input interface {
	// Get the current demand for this type of container
	GetDemand(containerType string) (int, error)
}
