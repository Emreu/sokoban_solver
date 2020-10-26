package world

import (
	"strings"
	"testing"
)

var metricLevel1 = `
#######
#@ $ .#
#######
`

func TestMetricSimple(t *testing.T) {
	r := strings.NewReader(metricLevel1)
	m, err := ReadMap(r)
	if err != nil {
		t.Fatalf("Level reading error: %v", err)
	}
	mc := NewMetricCalculator(m, Bitmap{})
	goalPos := Pos{5, 1}
	// sample dist at some points
	for _, x := range []int{2, 3, 4, 5} {
		metrics := mc.field[1][x]
		if len(metrics) != 1 {
			t.Fatalf("Wrong number of metrics targets @(%d, 1): %v", x, metrics)
		}

		expected := 5 - x
		if metrics[0] != expected {
			t.Errorf("Wrong metric value to goal @%s at @(%d,1): expected %d, found %d", goalPos, x, expected, metrics[0])
		}
	}
	metrics := mc.field[1][1]
	if len(metrics) != 1 {
		t.Fatalf("Wrong number of metrics targets @(1, 1): %v", metrics)
	}
	if metrics[0] != -1 {
		t.Fatalf("Wrong value of metrics targets @(1, 1): %v", metrics[0])
	}
}

var metricLevel2 = `
######
#@$ .#
# $$.#
#  ###
#. #
####
`

func TestMetricEvaluation(t *testing.T) {
	r := strings.NewReader(metricLevel2)
	m, err := ReadMap(r)
	if err != nil {
		t.Fatalf("Level reading error: %v", err)
	}
	mc := NewMetricCalculator(m, Bitmap{})

	state := State{
		MoveDomain:   NewMoveDomain(),
		BoxPositions: m.InitialBoxPositions(),
	}

	metric, err := mc.Evaluate(state)
	if err != nil {
		t.Fatal("Valid state has been evaluated as failed")
	}
	expected := 6
	if metric != expected {
		t.Fatalf("Wrong metric value (initial state) expected: %d got: %d", expected, metric)
	}

	state.BoxPositions = []Pos{
		{3, 1},
		{3, 2},
		{2, 3},
	}

	metric, err = mc.Evaluate(state)
	if err != nil {
		t.Fatal("Valid state has been evaluated as failed")
	}
	expected = 6
	if metric != expected {
		t.Fatalf("Wrong metric value expected: %d got: %d", expected, metric)
	}

	state.BoxPositions = []Pos{
		{2, 1},
		{3, 1},
		{3, 2},
	}

	metric, err = mc.Evaluate(state)
	if err == nil {
		t.Fatalf("Failed state has been evaluated as valid")
	}
}
