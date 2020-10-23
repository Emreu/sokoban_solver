package world

import "encoding/json"

type BoxMove struct {
	BoxIndex  int
	Direction MoveDirection
}

type NodeFail int

const (
	NodeOK NodeFail = iota
	NodeDuplicate
	NodeInvalid
)

type Node struct {
	ID     int64
	Metric int
	Hash   uint64
	Fail   NodeFail
	State
	// Parent node for easy backward traverse
	Parent *Node
	// Moves represent possible moves from current state with pointer to corresponding states
	Moves map[BoxMove]*Node
}

func NewNode(s State) *Node {
	return &Node{
		Metric: -1,
		Hash:   s.Hash(),
		State:  s,
		Moves:  make(map[BoxMove]*Node),
	}
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
	N.Boxes = n.State.BoxPositions
	N.Domain = n.State.MoveDomain.ListPosition()
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
