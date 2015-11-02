// The demand package defines the interface for demand models
package demand

type Task struct {
	Demand     int
	Requested  int
	Running    int
	FamilyName string
	Image      string
	Command    string
}

type Input interface {
	// Get the current demand for this type of container
	GetDemand(containerType string) (int, error)
}
