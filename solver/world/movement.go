package world

type MoveDirection int

const (
	MoveUp MoveDirection = iota
	MoveRight
	MoveDown
	MoveLeft
)

type MoveDomain struct{}

func NewMoveDomain() MoveDomain {
	return MoveDomain{}
}

func (md *MoveDomain) HasPosition(pos Position) bool {
	// TODO: implement
	return false
}
