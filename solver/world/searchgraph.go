package world

import (
	"encoding/json"
	"fmt"
)

type BoxMove struct {
	BoxIndex  int
	Distance  int
	Direction MoveDirection
}

func (b BoxMove) String() string {
	return fmt.Sprintf("[Move #%d %s (%d)]", b.BoxIndex, b.Direction, b.Distance)
}

type NodeFail int

const (
	NodeOK NodeFail = iota
	NodeDuplicate
	NodeInvalid
	NodeDeadlock
)

type Node struct {
	ID         int64
	Metric     int
	Hash       uint64
	Fail       NodeFail
	MoveDomain Bitmap
	Boxes      PosList
	// Parent node for easy backward traverse
	Parent *Node
	// Moves represent possible moves from current state with pointer to corresponding states
	Moves map[BoxMove]*Node
}

func (n Node) ApplyMove(move BoxMove) PosList {
	newBoxes := make(PosList, len(n.Boxes))
	copy(newBoxes, n.Boxes)
	newBoxes[move.BoxIndex] = newBoxes[move.BoxIndex].MoveInDirection(move.Direction)
	return newBoxes
}

func (n Node) MarshalJSON() ([]byte, error) {
	var N struct {
		ID     int64  `json:"id"`
		Parent int64  `json:"parent"`
		Metric int    `json:"metric"`
		Hash   uint64 `json:"hash"`
		Boxes  []Pos  `json:"boxes"`
		Domain []Pos  `json:"domain"`
		Fail   string `json:"fail,omitempty"`
	}
	N.ID = n.ID
	if n.Parent != nil {
		N.Parent = n.Parent.ID
	}
	N.Metric = n.Metric
	N.Hash = n.Hash
	N.Boxes = n.Boxes
	N.Domain = n.MoveDomain.List()
	switch n.Fail {
	case NodeDuplicate:
		N.Fail = "duplicate"
	case NodeInvalid:
		N.Fail = "invalid"
	}

	return json.Marshal(N)
}

// DeadEndNode represent constant reference for termination of dead end paths in search graphs
var DeadEndNode = &Node{}
