package world

import (
	"strings"
	"testing"
)

var basicLevel = `
#####
#@$.#
#####
`

func TestMapLoad(t *testing.T) {
	r := strings.NewReader(basicLevel)
	m, err := ReadMap(r)
	if err != nil {
		t.Fatalf("Simple map reading error: %v", err)
	}
	if m.Width != 5 {
		t.Errorf("unexpected map width: %d", m.Width)
	}
	if m.Height != 3 {
		t.Errorf("unexpected map height: %d", m.Height)
	}
	if m.StartX != 1 || m.StartY != 1 {
		t.Errorf("incorrect start position: (%d,%d)", m.StartX, m.StartY)
	}
}
