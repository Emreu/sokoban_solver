package world

type BoxMove struct {
	BoxIndex  int
	Direction MoveDirection
}

type SearchGraph struct {
	Root *Node
	// TODO: add state index for fast checking
}

type Node struct {
	State
	// Parent node for easy backward traverse
	Parent *Node
	// Moves represent possible moves from current state with pointer to corresponding states
	Moves map[BoxMove]*Node
}

// DeadEndNode represent constant reference for termination of dead end paths in search graphs
var DeadEndNode = &Node{}
