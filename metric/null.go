package metric

// Metric for cases such as Remainder rules, where we don't need to actually measure a current value
type NullMetric struct{}

// compile-time assert that we implement the right interface
var _ Metric = (*NullMetric)(nil)

func NewNullMetric() *NullMetric {
	return &NullMetric{}
}

func (n *NullMetric) UpdateCurrent() {}

func (n *NullMetric) Current() int {
	return 0
}
