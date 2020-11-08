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

/*
var exploreLevel = `
######
#@$ .#
# $$.#
#  ###
#. #
####
`

func TestExplorationLogic(t *testing.T) {
	r := strings.NewReader(exploreLevel)
	m, err := ReadMap(r)
	if err != nil {
		t.Fatalf("Level reading error: %v", err)
	}

	solver := NewSolver(m)

	solver.deadZones = FindDeadZones(m)

	// solver.metricCalc = NewMetricCalculator(solver.Map, solver.deadZones)

	domain := NewMoveDomain()
	domain.AddPosition(Pos{1, 1})
	domain.AddPosition(Pos{2, 1})
	domain.AddPosition(Pos{3, 1})
	domain.AddPosition(Pos{1, 2})
	domain.AddPosition(Pos{2, 2})
	domain.AddPosition(Pos{1, 3})
	domain.AddPosition(Pos{1, 4})
	domain.AddPosition(Pos{2, 4})

	boxPositions := []Pos{
		{4, 1},
		{3, 2},
		{2, 3},
	}

	newDomain, moves := NewMoveDomainFromMap(m, boxPositions, Pos{2, 2}, solver.deadZones)

	t.Logf("New domain: %s", newDomain)
	t.Fatalf("Moves: %+v", moves)

}
*/
