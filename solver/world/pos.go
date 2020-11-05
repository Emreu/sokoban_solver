package world

import (
	"fmt"
	"hash/fnv"
	"sort"
)

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

type PosList []Pos

func (p PosList) Hash() uint64 {
	hash := fnv.New64()

	for _, pos := range p.Sorted() {
		hash.Write([]byte{byte(pos.X)})
		hash.Write([]byte{byte(pos.Y)})
	}

	return hash.Sum64()
}

func (p PosList) Sorted() PosList {
	sort.Slice(p, func(i, j int) bool {
		if p[i].Y < p[j].Y {
			return true
		} else if p[i].Y > p[j].Y {
			return false
		} else if p[i].X < p[j].X {
			return true
		}
		return false
	})
	return p
}

func SortedPositions(pos []Pos) []Pos {
	sort.Slice(pos, func(i, j int) bool {
		if pos[i].Y < pos[j].Y {
			return true
		} else if pos[i].Y > pos[j].Y {
			return false
		} else if pos[i].X < pos[j].X {
			return true
		}
		return false
	})
	return pos
}
