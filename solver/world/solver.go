package world

import (
	"context"
	"fmt"
	"log"
)

type Solver struct {
	Map
	deadZones  MoveDomain
	root       *Node
	solution   *Node
	metricCalc MetricCalculator
	Done       bool
}

func NewSolver(m Map) *Solver {
	return &Solver{
		Map:       m,
		deadZones: NewMoveDomain(),
	}
}

// Find dead zones - positions from where it is impossible to move box to goal
func (s *Solver) findDeadZones() {
	var deadCorners []Pos
	addDeadZone := func(x, y int) {
		pos := Pos{x, y}
		if !s.Map.IsInside(pos) {
			return
		}
		switch s.Map.AtPos(pos) {
		case TileWall, TileGoal, TilePlayerOnGoal, TileBoxOnGoal:
			// do nothing
			return
		}
		// don't add if already added
		if s.deadZones.HasPosition(pos) {
			return
		}
		deadCorners = append(deadCorners, pos)
		s.deadZones.AddPosition(pos)
	}
	// at rirst find dead corners
	for y, row := range s.Map.Tiles {
		for x, tile := range row {
			if tile == TileWall {
				var up, down, left, right bool
				// check diagonal positions
				if s.Map.At(x-1, y-1) == TileWall {
					up = true
					left = true
				}
				if s.Map.At(x+1, y-1) == TileWall {
					up = true
					right = true
				}
				if s.Map.At(x+1, y+1) == TileWall {
					right = true
					down = true
				}
				if s.Map.At(x-1, y+1) == TileWall {
					left = true
					down = true
				}
				// add corresponding postition to dead zones if there is no wall or goal
				if up {
					addDeadZone(x, y-1)
				}
				if right {
					addDeadZone(x+1, y)
				}
				if down {
					addDeadZone(x, y+1)
				}
				if left {
					addDeadZone(x-1, y)
				}
			}
		}
	}

	log.Printf("Deadzones (corners only): %s", s.deadZones)

	// propagate dead corners
	for _, startPos := range deadCorners {
		for _, dir := range []MoveDirection{MoveUp, MoveRight, MoveDown, MoveLeft} {
			// log.Printf("Propagating from %s to %s", startPos, dir)
			pos := startPos.MoveInDirection(dir)
			if s.Map.AtPos(pos) == TileWall {
				// log.Print("Stop: wall at first tile")
				continue
			}
			var leftOpen, rightOpen bool
			var deadCorridor []Pos
			// advance position & check left & right walls
		rayScan:
			for {
				// if out of map - stop
				if !s.Map.IsInside(pos) {
					// log.Print("Stop: out of map")
					break rayScan
				}
				// if goal on path - it not dead zone
				switch s.Map.AtPos(pos) {
				case TileGoal, TileBoxOnGoal, TilePlayerOnGoal:
					// log.Print("Stop: goal on line")
					break rayScan
				}
				// check left if it's not open yet
				if !leftOpen {
					leftPos := pos.MoveInDirection(dir.RotateCCW())
					if s.Map.AtPos(leftPos) != TileWall {
						leftOpen = true
					}
				}
				// check right if it's not open yet
				if !rightOpen {
					rightPos := pos.MoveInDirection(dir.RotateCW())
					if s.Map.AtPos(rightPos) != TileWall {
						rightOpen = true
					}
				}
				// if both open - stop propagation
				if leftOpen && rightOpen {
					// log.Print("Stop: both sides clear")
					break rayScan
				}
				// move ahead
				deadCorridor = append(deadCorridor, pos)
				pos = pos.MoveInDirection(dir)
				// if wall hit - stop propagation and add corridor to dead zone
				if s.Map.AtPos(pos) == TileWall {
					// log.Print("Wall hit - adding to deadzones!")
					for _, p := range deadCorridor {
						s.deadZones.AddPosition(p)
					}
					break rayScan
				}
			}
		}
	}

	log.Printf("Deadzones (with propagation): %s", s.deadZones)
}

func (s *Solver) Solve(c context.Context, debugOnly bool) error {
	if s.Done {
		return nil
	}
	// prepare
	log.Print("Finding dead zones...")
	s.findDeadZones()

	log.Print("Preparing metric calc...")
	s.metricCalc = NewMetricCalculator(s.Map, s.deadZones)

	if debugOnly {
		return nil
	}

	// initialize
	log.Print("Initializing...")
	boxPositions := s.Map.InitialBoxPositions()
	state := State{
		MoveDomain:   NewMoveDomainFromMap(s.Map, boxPositions, s.Map.StartPos()),
		BoxPositions: boxPositions,
	}
	root := NewNode(state)
	if !s.populateMoves(root) {
		return fmt.Errorf("no moves available on init")
	}
	s.root = root

	if s.isSolution(state) {
		return nil
	}

	// do breadth first search
	var exploreFrontier = []*Node{root}

	step := 0
	for len(exploreFrontier) != 0 {
		var nextFrontier []*Node
		log.Printf("Exploring frontier #%d of %d nodes...", step, len(exploreFrontier))
		for _, node := range exploreFrontier {
			newNodes := s.exploreNode(node)
			for _, n := range newNodes {
				// check if solution was found
				if s.isSolution(n.State) {
					log.Print("Solution found!")
					s.solution = n
					s.Done = true
					return nil
				}
				nextFrontier = append(nextFrontier, n)
			}
		}
		exploreFrontier = nextFrontier
		// TODO: add explore frontier limitation if its too large
		// range nodes based on state estimation and explore only top-N nodes
		// store best metric value for postponed nodes and reconsider them only if current frontier is too small or metrics of its nodes become worse
		step++
	}

	s.Done = true
	return nil
}

func (s Solver) exploreNode(n *Node) []*Node {
	// log.Printf("Exploring node %p:", n)
	var nextNodes []*Node
	for move := range n.Moves {
		// log.Printf("+ move #%d %s", move.BoxIndex, move.Direction)
		// copy box positions from current state
		var nextPositions = make([]Pos, len(n.State.BoxPositions))
		copy(nextPositions, n.State.BoxPositions)
		// apply move to current box positions
		boxSrcPos := nextPositions[move.BoxIndex]
		boxDstPos := boxSrcPos.MoveInDirection(move.Direction)
		nextPositions[move.BoxIndex] = boxDstPos

		state := State{
			MoveDomain:   NewMoveDomainFromMap(s.Map, nextPositions, boxSrcPos),
			BoxPositions: nextPositions,
		}

		// TODO: check if state is already indexed

		// create node
		node := NewNode(state)
		node.Parent = n
		s.populateMoves(node)
		n.Moves[move] = node
		nextNodes = append(nextNodes, node)
	}

	return nextNodes
}

// populateMoves find available box movements from node state and adds them to moves map
// return false if no moves available
func (s Solver) populateMoves(n *Node) bool {
	// log.Print("Populating moves...")
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
			// check destination isn't in dead zone
			if s.deadZones.HasPosition(dstPos) {
				continue
			}

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
	current := s.solution
	previous := current.Parent
	for previous != nil {
		for move := range previous.Moves {
			if previous.Moves[move] != current {
				continue
			}
			log.Printf("Move #%d %s", move.BoxIndex, move.Direction)
		}
		current = previous
		previous = current.Parent
	}
	// TODO: use A* to generate path from state to state (in move domain only)
	// reverse path
	return []MoveDirection{MoveUp, MoveDown, MoveLeft, MoveRight}, nil
}

type SolverDebug struct {
	DeadZones []Pos            `json:"dead_zones"`
	Metrics   []map[string]int `json:"metrics"`
}

func (s Solver) GetDebug() SolverDebug {
	dz := s.deadZones.ListPosition()
	var metricsMap = make(map[Pos]map[string]int)

	for y, row := range s.metricCalc.cells {
		for x, cell := range row {
			for goalPos, value := range cell.dist {
				field, ok := metricsMap[goalPos]
				if !ok {
					field = make(map[string]int)
				}
				field[fmt.Sprintf("%d,%d", x, y)] = value
				metricsMap[goalPos] = field
			}
		}
	}

	var metricList []map[string]int
	for _, field := range metricsMap {
		metricList = append(metricList, field)
	}

	return SolverDebug{
		DeadZones: dz,
		Metrics:   metricList,
	}
}
