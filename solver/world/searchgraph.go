package world

type BoxMove struct {
	BoxIndex  int
	Direction MoveDirection
}

type Node struct {
	State
	// Parent node for easy backward traverse
	Parent *Node
	// Moves represent possible moves from current state with pointer to corresponding states
	Moves map[BoxMove]*Node
}

func NewNode(s State) *Node {
	return &Node{
		State: s,
		Moves: make(map[BoxMove]*Node),
	}
}

// DeadEndNode represent constant reference for termination of dead end paths in search graphs
var DeadEndNode = &Node{}
