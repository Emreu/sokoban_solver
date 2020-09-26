package world

import (
	"context"
	"fmt"
)

type Solver struct {
	Map
	SearchGraph
	Done bool
}

func NewSolver(m Map) *Solver {
	return &Solver{
		Map:         m,
		SearchGraph: SearchGraph{},
	}
}

func (s *Solver) Solve(c context.Context) error {
	if s.Done {
		return nil
	}
	// TODO: implement
	// for every node in search graph
	// find unexplored movements
	// for every movement generate new state (if not registered before)
	// add nodes for corresponding moves and states
	// evaluate metrics for new states
	// range states by metrics

	s.Done = true
	return nil
}

func (s Solver) GetPath() ([]MoveDirection, error) {
	if !s.Done {
		return nil, fmt.Errorf("call Solve() before GetPath()")
	}
	// TODO: traverse search graph from solution node back to root
	return []MoveDirection{MoveUp, MoveDown, MoveLeft, MoveRight}, nil
}
