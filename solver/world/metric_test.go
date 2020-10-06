package world

import (
	"strings"
	"testing"
)

var metricLevel = `
#######
#. + *# 
#######
`

func TestMetricSimple(t *testing.T) {
	r := strings.NewReader(metricLevel)
	m, err := ReadMap(r)
	if err != nil {
		t.Fatalf("Level reading error: %v", err)
	}
	mc := NewMetricCalculator(m)
	// sample dist at some points
	_ = mc
	// TODO: write tests
}
