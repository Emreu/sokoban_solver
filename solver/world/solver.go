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

	// initialize
	boxPositions := s.Map.InitialBoxPositions()
	state := State{
		MoveDomain:   NewMoveDomainFromMap(s.Map, boxPositions, s.Map.StartPos()),
		BoxPositions: boxPositions,
	}
	root := &Node{
		State: state,
	}
	if !s.populateMoves(root) {
		return fmt.Errorf("no moves available on init")
	}
	s.SearchGraph.Root = root

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

// populateMoves find available box movements from node state and adds them to moves map
// return false if no moves available
func (s Solver) populateMoves(n *Node) bool {
	movesFound := false
	// loop over all boxes on map
	// TODO: consider optimization with contact positions storing
	for i, boxPos := range n.State.BoxPositions {
		// find valid moves
		directions := []MoveDirection{MoveUp, MoveRight, MoveDown, MoveLeft}
	moves:
		for _, dir := range directions {
			// move is valid if its done from move domain and final tile is empty and not in forbidden zone
			srcPos := boxPos.MoveAgainstDirection(dir)
			if !n.State.MoveDomain.HasPosition(srcPos) {
				continue
			}
			dstPos := boxPos.MoveInDirection(dir)
			if s.Map.AtPos(dstPos) == TileWall {
				continue
			}
			// check if any box occupies destination tile
			for _, p := range n.State.BoxPositions {
				if p == dstPos {
					continue moves
				}
			}
			// TODO: add forbidden zones checking

			// finally add move
			n.Moves[BoxMove{
				BoxIndex:  i,
				Direction: dir,
			}] = nil
			movesFound = true
		}
	}
	return movesFound
}

func (s Solver) isSolution(state State) bool {
	for _, pos := range state.BoxPositions {
		switch s.Map.AtPos(pos) {
		case TileBoxOnGoal, TilePlayerOnGoal, TileGoal:
			continue
		default:
			return false
		}
	}
	return true
}

func (s Solver) GetPath() ([]MoveDirection, error) {
	if !s.Done {
		return nil, fmt.Errorf("call Solve() before GetPath()")
	}
	// TODO: traverse search graph from solution node back to root
	return []MoveDirection{MoveUp, MoveDown, MoveLeft, MoveRight}, nil
}
