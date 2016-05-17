package metric

// ToyMetric is only used for testing, but we can set its value
type ToyMetric struct {
	SettableCurrent int
}

// compile-time assert that we implement the right interface
var _ Metric = (*ToyMetric)(nil)

// NewToyMetric creates a new toy metric
func NewToyMetric() *ToyMetric {
	return &ToyMetric{}
}

// UpdateCurrent reads the value of the current metric, but this is a no-op for the Toy metric
func (t *ToyMetric) UpdateCurrent() {}

// Current reads out the value of the current queue length
func (t *ToyMetric) Current() int {
	return t.SettableCurrent
}
