package world

import (
	"fmt"
	"strings"
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

type MoveDomain struct {
	// tiles map[Pos]struct{} // simple implementation
	bitmap Bitmap
}

func NewMoveDomain() MoveDomain {
	return MoveDomain{
		// tiles: make(map[Pos]struct{}),
		bitmap: Bitmap{},
	}
}

// NewMoveDomainFromMap generate move domain from Map (considering only walls), box positions and starting position
// TODO: refactor to create state with predefined moves
func NewMoveDomainFromMap(m Map, boxPositions []Pos, start Pos, deadZones Bitmap) (MoveDomain, []BoxMove) {
	domain := NewMoveDomain()
	boxes := Bitmap{}
	var moves []BoxMove

	for _, pos := range boxPositions {
		boxes.SetBit(pos)
	}

	boxIndex := func(p Pos) int {
		for i, pos := range boxPositions {
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
		// if we hit a box - check if move is possible and save
		if boxes.CheckBit(p) {
			dstPos := p.MoveInDirection(d)
			// check if destination position is available
			if m.AtPos(dstPos) == TileWall {
				return
			}
			if boxes.CheckBit(dstPos) {
				return
			}
			// check destination isn't in dead zone
			if deadZones.CheckBit(dstPos) {
				return
			}

			moves = append(moves, BoxMove{
				Direction: d,
				Distance:  dist,
				BoxIndex:  boxIndex(p),
			})
			return
		}
		if domain.HasPosition(p) {
			return
		}
		domain.AddPosition(p)
		backward := d.Opposite()
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

func (md *MoveDomain) HasPosition(pos Pos) bool {
	return md.bitmap.CheckBit(pos)
}

func (md *MoveDomain) AddPosition(pos Pos) {
	md.bitmap.SetBit(pos)
}

func (md *MoveDomain) ListPosition() []Pos {
	return md.bitmap.List()
}

func (md MoveDomain) HashBytes() []byte {
	return md.bitmap.HashBytes()
}

func (md MoveDomain) String() string {
	buf := &strings.Builder{}
	buf.WriteString("MoveDomain{")
	for _, pos := range md.bitmap.List() {
		buf.WriteString(pos.String())
	}
	buf.WriteString("}")
	return buf.String()
}
