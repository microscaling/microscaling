package metric

// Metric for cases such as Remainder rules, where we don't need to actually measure a current value
type NullMetric struct{}

// compile-time assert that we implement the right interface
var _ Metric = (*NullMetric)(nil)

// NewNullMetric creates a new Null metric
func NewNullMetric() *NullMetric {
	return &NullMetric{}
}

// UpdateCurrent reads the value of the current metric, but this is a no-op for the Null metric
func (n *NullMetric) UpdateCurrent() {}

// Current reads out the value of the current queue length - which is always 0 for the Null metric
func (n *NullMetric) Current() int {
	return 0
}
