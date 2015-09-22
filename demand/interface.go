package demand

// Get the current demand for this type of container
type Input interface {
  GetDemand(containerType string) (int, error)
}