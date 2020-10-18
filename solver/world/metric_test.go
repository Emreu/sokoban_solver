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
	mc := NewMetricCalculator(m, NewMoveDomain())
	goalPos := Pos{5, 1}
	// sample dist at some points
	for _, x := range []int{2, 3, 4, 5} {
		d := mc.cells[1][x].dist
		if len(d) != 1 {
			t.Fatalf("Wrong number of metrics targets @(%d, 1): %v", x, d)
		}

		dist, ok := d[goalPos]
		if !ok {
			t.Errorf("Metric to goal @%s not found @(%d,1)", goalPos, x)
			continue
		}
		expected := 5 - x
		if dist != expected {
			t.Errorf("Wrong metric value to goal @%s at @(%d,1): expected %d, found %d", goalPos, x, expected, dist)
		}
	}
	d := mc.cells[1][1].dist
	if len(d) != 0 {
		t.Fatalf("Wrong number of metrics targets @(1, 1): %v", d)
	}
}
