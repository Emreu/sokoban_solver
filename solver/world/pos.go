package world

import "fmt"

type Pos struct {
	X int
	Y int
}

func (p Pos) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

func (p Pos) Neighbours() []Pos {
	return []Pos{
		{X: p.X, Y: p.Y - 1},
		{X: p.X + 1, Y: p.Y},
		{X: p.X, Y: p.Y + 1},
		{X: p.X - 1, Y: p.Y},
	}
}

func (p Pos) Diagonals() []Pos {
	return []Pos{
		{X: p.X - 1, Y: p.Y - 1},
		{X: p.X + 1, Y: p.Y - 1},
		{X: p.X + 1, Y: p.Y + 1},
		{X: p.X - 1, Y: p.Y + 1},
	}
}

func (p Pos) MoveInDirection(d MoveDirection) Pos {
	switch d {
	case MoveUp:
		return Pos{X: p.X, Y: p.Y - 1}
	case MoveRight:
		return Pos{X: p.X + 1, Y: p.Y}
	case MoveDown:
		return Pos{X: p.X, Y: p.Y + 1}
	case MoveLeft:
		return Pos{X: p.X - 1, Y: p.Y}
	}
	return p
}

func (p Pos) MoveAgainstDirection(d MoveDirection) Pos {
	switch d {
	case MoveUp:
		return Pos{X: p.X, Y: p.Y + 1}
	case MoveRight:
		return Pos{X: p.X - 1, Y: p.Y}
	case MoveDown:
		return Pos{X: p.X, Y: p.Y - 1}
	case MoveLeft:
		return Pos{X: p.X + 1, Y: p.Y}
	}
	return p
}

func (p Pos) Mirror(center Pos) Pos {
	return Pos{
		X: 2*center.X - p.X,
		Y: 2*center.Y - p.Y,
	}
}
