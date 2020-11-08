package world

import (
	"strings"
	"testing"
)

var exploreLevel = `
########  
#  .#  #  
#  .#  ###
#  ##    #
## $ $@  #
 # ##  ###
 #    ##  
 ######   
`

func TestExplore(t *testing.T) {
	r := strings.NewReader(exploreLevel)
	m, err := ReadMap(r)
	if err != nil {
		t.Fatalf("Simple map reading error: %v", err)
	}

	deadzones := FindDeadZones(m)
	hlock, vlock := FindMoveLocks(m, deadzones)

	domain, moves := Explore(m, m.InitialBoxPositions(), m.StartPos(), deadzones, hlock, vlock)

	t.Logf("Move domain: %s", domain)

	t.Errorf("Moves: %v", moves)

	if len(moves) != 3 {

		t.Fatalf("Unexpected number of moves: %d", len(moves))
	}
}
