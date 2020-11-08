package world

import (
	"fmt"
)

type MoveDirection int

const (
	MoveUp MoveDirection = iota
	MoveRight
	MoveDown
	MoveLeft
)

var AllDirections = []MoveDirection{MoveUp, MoveRight, MoveDown, MoveLeft}

func (d MoveDirection) String() string {
	switch d {
	case MoveUp:
		return "up"
	case MoveRight:
		return "right"
	case MoveDown:
		return "down"
	case MoveLeft:
		return "left"
	}
	return "???"
}

func (d MoveDirection) Opposite() MoveDirection {
	switch d {
	case MoveUp:
		return MoveDown
	case MoveRight:
		return MoveLeft
	case MoveDown:
		return MoveUp
	case MoveLeft:
		return MoveRight
	}
	return -1
}

func (d MoveDirection) MarshalJSON() ([]byte, error) {
	switch d {
	case MoveUp:
		return []byte(`"^"`), nil
	case MoveRight:
		return []byte(`">"`), nil
	case MoveLeft:
		return []byte(`"<"`), nil
	case MoveDown:
		return []byte(`"V"`), nil
	}
	return nil, fmt.Errorf("unknown direction: %d", d)
}

func (d MoveDirection) RotateCW() MoveDirection {
	d++
	if d > MoveLeft {
		d = MoveUp
	}
	return d
}

func (d MoveDirection) RotateCCW() MoveDirection {
	d--
	if d < MoveUp {
		d = MoveLeft
	}
	return d
}
