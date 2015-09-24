// The demand package defined the interface for demand models
package demand

type Input interface {
	// Get the current demand for this type of container
	GetDemand(containerType string) (int, error)
}
