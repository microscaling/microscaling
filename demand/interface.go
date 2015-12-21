// The demand package defines the interface for demand models
package demand

type Task struct {
	Demand          int
	Requested       int
	Running         int
	FamilyName      string
	Image           string
	Command         string
	PublishAllPorts bool
}
