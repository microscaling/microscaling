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
	// Update the demand for all the apps we're dealing with. Returns True if demand for anything changed
	Update(ts map[string]Task) (bool, error)
}
