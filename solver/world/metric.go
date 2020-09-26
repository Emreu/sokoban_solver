package world

type MetricCalculator struct{}

func NewMetricCalculator(m Map) MetricCalculator {
	// TODO: calculate distance fields for every goal
	return MetricCalculator{}
}

func (m MetricCalculator) Evaluate(s State) int {
	// TODO: implement
	return -1
}
