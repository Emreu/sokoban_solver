package world

import (
	"context"
	"errors"
	"fmt"
	"log"
)

var (
	errSolutionFound  = errors.New("solution found")
	errDuplicateState = errors.New("duplicate state")
	errDeadlock       = errors.New("deadlock")
	errInvalidMetric  = errors.New("invalid metric")
)

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
	s.hLock, s.vLock = FindMoveLocks(s.m, s.deadZones)

	log.Print("Preparing metric calc...")
	s.metricCalc = NewMetricCalculator(s.m, s.deadZones)

	// initialize
	log.Print("Starting solution search...")
	var err error
	s.root, err = s.createNode(s.m.InitialBoxPositions(), s.m.StartPos(), nil)
	if err == errSolutionFound {
		s.solution = s.root
		return nil
	}

	s.frontier.Add(s.root)

	for s.frontier.HasNodes() {
		for _, node := range s.frontier.GetBatch() {
			select {
			case <-c.Done():
				return fmt.Errorf("timed out or canceled")
			default:
			}
			for move := range node.Moves {
				playePos := node.Boxes[move.BoxIndex]
				newBoxes := node.ApplyMove(move)
				child, err := s.createNode(newBoxes, playePos, node)
				switch err {
				case nil:
					node.Moves[move] = child
					s.frontier.Add(child)
				case errSolutionFound:
					s.solution = child
					node.Moves[move] = child
					return nil
				default:
				}
			}
		}
	}

	return fmt.Errorf("all states explored but solution not found")
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
	DeadZones PosList          `json:"dead_zones"`
	HLock     PosList          `json:"h_locks"`
	VLock     PosList          `json:"v_locks"`
	Metrics   []map[string]int `json:"metrics"`
}

func (s Solver) GetDebug() SolverDebug {
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
		DeadZones: s.deadZones.List(),
		HLock:     s.hLock.List(),
		VLock:     s.vLock.List(),
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

func (s *Solver) createNode(boxes PosList, start Pos, parent *Node) (node *Node, err error) {
	// check if state was already tracked
	hash := boxes.Hash()
	hashNodes, hashExists := s.boxHashIndex[hash]
	if hashExists {
		// check all nodes with same box positions
		for _, n := range hashNodes {
			// if start is in move domain - state is the same
			if n.MoveDomain.CheckBit(start) {
				log.Printf("Duplicate node skipped!")
				return n, errDuplicateState
			}
		}
	}
	// create node
	node = &Node{
		ID:     s.getID(),
		Parent: parent,
		Moves:  make(map[BoxMove]*Node),
		Boxes:  boxes,
	}
	// defer save hash - valid or dead end
	defer func() {
		if err != nil {
			node.MoveDomain = ExploreDomainOnly(s.m, boxes, start)
		}
		hashNodes = append(hashNodes, node)
		s.boxHashIndex[hash] = hashNodes
	}()
	// check if state is solution, but only if parent is none (initial state) or parent metric = 1 (one step to solution)
	if parent == nil || parent.Metric == 1 {
		if s.isSolution(boxes) {
			return node, errSolutionFound
		}
	}
	// TODO: check if deadlock detection is required

	// build move domain and box moves
	domain, moves := Explore(s.m, boxes, start, s.deadZones, s.hLock, s.vLock)
	node.MoveDomain = domain
	for _, move := range moves {
		node.Moves[move] = nil
	}

	// evaluate metric
	metric, err := s.metricCalc.Evaluate(boxes)
	if err != nil {
		node.Fail = NodeInvalid
		return node, errInvalidMetric
	}
	node.Metric = metric

	return node, nil
}
