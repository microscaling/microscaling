package demand

type Input interface {
	// Get the current demand for this type of container
	GetDemand(containerType string) (int, error)
}
