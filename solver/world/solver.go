package world

import (
	"context"
	"fmt"
	"log"
)

const ExploreBatch = 200

type Solver struct {
	m             Map
	deadZones     Bitmap
	hLock         Bitmap
	vLock         Bitmap
	root          *Node
	solution      *Node
	nodeHashIndex map[uint64]struct{}
	boxHashIndex  map[uint64][]*Node
	metricCalc    MetricCalculator
	nextID        int64
	frontier      Frontier
}

func NewSolver(m Map) *Solver {
	return &Solver{
		m:             m,
		deadZones:     Bitmap{},
		nodeHashIndex: make(map[uint64]struct{}),
		boxHashIndex:  make(map[uint64][]*Node),
	}
}

func mergeSorted(base, addition []*Node) []*Node {
	var result = make([]*Node, len(base)+len(addition))
	var i, j, k = 0, 0, 0
	for i < len(base) && j < len(addition) {
		if base[i].Metric < addition[j].Metric {
			result[k] = base[i]
			i++
		} else {
			result[k] = addition[j]
			j++
		}
		k++
	}
	if i < len(base) {
		copy(result[k:], base[i:])
	} else {
		copy(result[k:], addition[j:])
	}
	return result
}

func (s *Solver) getID() int64 {
	s.nextID++
	return s.nextID
}

func (s *Solver) Solve(c context.Context) error {
	if s.solution != nil {
		return nil
	}
	// prepare
	log.Print("Finding dead zones...")
	s.deadZones = FindDeadZones(s.m)

	log.Print("Finding h/v lock zones...")
	s.hLock, s.vLock = FindMoveLocks(s.m)

	log.Print("Preparing metric calc...")
	s.metricCalc = NewMetricCalculator(s.m, s.deadZones)

	// initialize
	log.Print("Starting solution search...")
	s.root = s.createNode(s.m.InitialBoxPositions(), s.m.StartPos(), nil)
	// TODO: check

	s.frontier.Add(s.root)

	for {
		node, err := s.frontier.Pick()
		if err != nil {
			return fmt.Errorf("error picking next node: %v", err)
		}

	}

	return fmt.Errorf("all states explored but solution not found")
}

func (s *Solver) exploreNode(n *Node) []*Node {
	// log.Printf("Exploring node %p:", n)
	var nextNodes []*Node
	for move := range n.Moves {
		// log.Printf("+ move #%d %s", move.BoxIndex, move.Direction)
		// copy box positions from current state
		var nextPositions = make([]Pos, len(n.Boxes))
		copy(nextPositions, n.Boxes)
		// apply move to current box positions
		boxSrcPos := nextPositions[move.BoxIndex]
		boxDstPos := boxSrcPos.MoveInDirection(move.Direction)
		nextPositions[move.BoxIndex] = boxDstPos

		domain, moves := NewMoveDomainFromMap(s.m, nextPositions, boxSrcPos, s.deadZones)

		state := State{
			MoveDomain:   domain,
			BoxPositions: nextPositions,
		}

		// create node
		node := NewNode(state)
		node.ID = s.getID()
		node.Parent = n

		// check if state isn't indexed
		if _, found := s.nodeHashIndex[node.Hash]; found {
			node.Fail = NodeDuplicate
			n.Moves[move] = node
			continue
		}
		s.nodeHashIndex[node.Hash] = struct{}{}

		// calculate metric
		var err error
		node.Metric, err = s.metricCalc.Evaluate(state.BoxPositions)
		if err != nil {
			node.Fail = NodeInvalid
			n.Moves[move] = node
			continue
		}

		setMoves(node, moves)
		n.Moves[move] = node

		nextNodes = append(nextNodes, node)
	}

	return nextNodes
}

func setMoves(n *Node, moves []BoxMove) bool {
	for _, m := range moves {
		n.Moves[m] = nil
	}
	return len(moves) > 0
}

func (s Solver) isSolution(boxes PosList) bool {
	for _, pos := range boxes {
		switch s.m.AtPos(pos) {
		case TileBoxOnGoal, TilePlayerOnGoal, TileGoal:
			continue
		default:
			return false
		}
	}
	return true
}

func (s Solver) GetPath() ([]MoveDirection, error) {
	if s.solution == nil {
		return nil, fmt.Errorf("call Solve() before GetPath()")
	}
	var path []MoveDirection
	current := s.solution
	previous := current.Parent
	var dest Pos
	lastState := true

	for previous != nil {
		var move BoxMove
		moveFound := false
		for m := range previous.Moves {
			if previous.Moves[m] == current {
				move = m
				moveFound = true
				break
			}
		}
		if !moveFound {
			return nil, fmt.Errorf("incorrect state tree!")
		}
		// log.Printf("Move #%d %s", move.BoxIndex, move.Direction)

		// start from where box was before move apply
		start := previous.Boxes[move.BoxIndex]
		if !lastState {
			// do transition to next destination if it's not last state
			var err error
			segment, err := FindDirections(current.MoveDomain, start, dest)
			if err != nil {
				return nil, fmt.Errorf("path finding error: %v", err)
			}
			if len(segment) > 0 {
				path = append(segment, path...)
			}
		}
		lastState = false

		// add move to path
		path = append([]MoveDirection{move.Direction}, path...)

		// save destination for next (previous ?_?) transition - tile before box
		dest = start.MoveAgainstDirection(move.Direction)

		// advance states
		current = previous
		previous = current.Parent
	}
	// finally find directions from start position to first box move
	segment, err := FindDirections(current.MoveDomain, s.m.StartPos(), dest)
	if err != nil {
		return nil, fmt.Errorf("path finding error: %v", err)
	}
	path = append(segment, path...)
	return path, nil
}

type SolverDebug struct {
	DeadZones []Pos            `json:"dead_zones"`
	Metrics   []map[string]int `json:"metrics"`
}

func (s Solver) GetDebug() SolverDebug {
	dz := s.deadZones.List()
	// var metricsMap = make(map[Pos]map[string]int)

	// for y, row := range s.metricCalc.field {
	// 	for x, metrics := range row {
	// 		for goalPos, value := range cell.dist {
	// 			field, ok := metricsMap[goalPos]
	// 			if !ok {
	// 				field = make(map[string]int)
	// 			}
	// 			field[fmt.Sprintf("%d,%d", x, y)] = value
	// 			metricsMap[goalPos] = field
	// 		}
	// 	}
	// }

	// var metricList []map[string]int
	// for _, field := range metricsMap {
	// 	metricList = append(metricList, field)
	// }

	return SolverDebug{
		DeadZones: dz,
		// Metrics:   metricList,
	}
}

func (s Solver) GetTree(max int) []*Node {
	var output []*Node

	var frontier = []*Node{s.root}

	for len(frontier) > 0 && (max == 0 || len(output) < max) {
		var nextFrontier []*Node
		for _, n := range frontier {
			output = append(output, n)
			for _, next := range n.Moves {
				if next == nil {
					continue
				}
				nextFrontier = append(nextFrontier, next)
			}
		}
		frontier = nextFrontier
	}

	return output
}

func (s *Solver) createNode(boxes PosList, start Pos, parent *Node) *Node {
	// check if state was already tracked
	hash := boxes.Hash()
	if nodes, exists := s.boxHashIndex[hash]; exists {
		// check all nodes with same box positions
		for _, n := range nodes {
			// if start is in move domain - state is the same
			if n.MoveDomain.CheckBit(start) {
				// TODO: return error for duplicate state
				return nil
			}
		}
		// TODO: save node by hash
	}
	// create node
	node := &Node{
		ID:     s.getID(),
		Parent: parent,
	}
	// TODO: save state to hash index
	// check if state is solution, but only if parent is none (initial state) or parent metric = 1 (one step to solution)
	if parent == nil || parent.Metric == 1 {
		if s.isSolution(boxes) {
			// TODO: return error for solution
			return nil
		}
	}
	// TODO: check if deadlock detection is required

	// build move domain and box moves
	domain, moves := BuildState(s.m, boxes, start, s.deadZones, s.hLock, s.vLock)
	node.MoveDomain = domain
	for _, move := range moves {
		node.Moves[move] = nil
	}

	// evaluate metric
	metric, err := s.metricCalc.Evaluate(boxes)
	if err != nil {
		// TODO: return invalid state error
		return nil
	}
	node.Metric = metric

	return node
}

func BuildState(m Map, boxes PosList, start Pos, deadZones, hLock, vLock Bitmap) (Bitmap, []BoxMove) {
	var domain Bitmap
	var boxBitmap Bitmap
	var moves []BoxMove

	for _, pos := range boxes {
		boxBitmap.SetBit(pos)
	}

	boxIndex := func(p Pos) int {
		for i, pos := range boxes {
			if p == pos {
				return i
			}
		}
		return -1
	}

	var fill func(p Pos, d MoveDirection, dist int)
	fill = func(p Pos, d MoveDirection, dist int) {
		if !m.IsInside(p) {
			return
		}
		if m.AtPos(p) == TileWall {
			return
		}
		backward := d.Opposite()
		// if we hit a box - check if move is possible and save
		if boxBitmap.CheckBit(p) {
			dstPos := p.MoveInDirection(d)
			// check if destination position is available
			if m.AtPos(dstPos) == TileWall {
				return
			}
			if boxBitmap.CheckBit(dstPos) {
				return
			}
			// check destination isn't in dead zone
			if deadZones.CheckBit(dstPos) {
				return
			}
			// for quick deadlock detection check if there is neighbour box with same h/v lock as target position
			dstHLock := hLock.CheckBit(dstPos)
			dstVLock := vLock.CheckBit(dstPos)
			if dstHLock || dstVLock {
				for _, dir := range AllDirections {
					if dir == backward {
						continue
					}
					n := p.MoveInDirection(dir)
					if !boxBitmap.CheckBit(n) {
						continue
					}
					if dstHLock && hLock.CheckBit(n) {
						return
					}
					if dstVLock && vLock.CheckBit(n) {
						return
					}
				}
			}
			// save move
			move := BoxMove{
				Direction: d,
				Distance:  dist,
				BoxIndex:  boxIndex(p),
			}
			moves = append(moves, move)
			return
		}
		if domain.CheckBit(p) {
			return
		}
		domain.SetBit(p)
		for _, dir := range AllDirections {
			if dir == backward {
				continue
			}
			n := p.MoveInDirection(dir)
			fill(n, dir, dist+1)
		}
	}

	fill(start, -1, 0)

	return domain, moves
}
