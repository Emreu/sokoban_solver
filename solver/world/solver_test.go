package world

import (
	"context"
	"strings"
	"testing"
)

var solverLevel = `
    #####
    #   #
    #$  #
  ###  $##
  #  $ $ #
### # ## #   ######
#   # ## #####  ..#
# $  $          ..#
##### ### #@##  ..#
    #     #########
    #######
`

func TestSolver(t *testing.T) {
	r := strings.NewReader(solverLevel)
	m, err := ReadMap(r)
	if err != nil {
		t.Fatalf("Level reading error: %v", err)
	}

	solver := NewSolver(m)

	ctx := context.Background()

	err = solver.Solve(ctx)
	if err != nil {
		t.Fatalf("Solving error: %v", err)
	}
}
