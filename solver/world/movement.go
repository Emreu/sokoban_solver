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
func NewMoveDomainFromMap(m Map, boxPositions []Pos, start Pos) MoveDomain {
	var perimeter = make(map[Pos]struct{})
	perimeter[start] = struct{}{}
	domain := NewMoveDomain()
	boxes := Bitmap{}

	for _, pos := range boxPositions {
		boxes.SetBit(pos)
	}

	var fill func(p Pos)
	fill = func(p Pos) {
		if !m.IsInside(p) {
			return
		}
		if m.AtPos(p) == TileWall {
			return
		}
		if boxes.CheckBit(p) {
			return
		}
		if domain.HasPosition(p) {
			return
		}
		domain.AddPosition(p)
		for _, n := range p.Neighbours() {
			fill(n)
		}
	}

	fill(start)

	// var nextPerimeter = make(map[Pos]struct{})
	// for {
	// 	for pos := range perimeter {
	// 		// skip if this tile is wall
	// 		if m.AtPos(pos) == TileWall {
	// 			continue
	// 		}
	// 		// skip if tile is occupied by box
	// 		if _, occupied := boxPos[pos]; occupied {
	// 			// but save contact position
	// 			// domain.AddContacts(pos)
	// 			continue
	// 		}
	// 		// add tile to domain
	// 		domain.AddPosition(pos)
	// 		// schedule neighbour tiles for next perimeter
	// 		for _, p := range pos.Neighbours() {
	// 			// skip if out of map
	// 			if !m.IsInside(p) {
	// 				continue
	// 			}
	// 			// skip if already processing
	// 			if _, processing := perimeter[p]; processing {
	// 				continue
	// 			}
	// 			// skip if already scheduled
	// 			if _, scheduled := nextPerimeter[p]; scheduled {
	// 				continue
	// 			}
	// 			// skip if already in domain
	// 			if domain.HasPosition(p) {
	// 				continue
	// 			}
	// 			// finally add to next perimeter
	// 			nextPerimeter[p] = struct{}{}
	// 		}
	// 	}
	// 	if len(nextPerimeter) == 0 {
	// 		break
	// 	}
	// 	perimeter = nextPerimeter
	// 	nextPerimeter = make(map[Pos]struct{})
	// }

	return domain
}

func (md *MoveDomain) HasPosition(pos Pos) bool {
	return md.bitmap.CheckBit(pos)
	// _, exists := md.tiles[pos]
	// return exists
}

func (md *MoveDomain) AddPosition(pos Pos) {
	md.bitmap.SetBit(pos)
	// md.tiles[pos] = struct{}{}
}

func (md *MoveDomain) ListPosition() []Pos {
	return md.bitmap.List()
	// var pos []Pos
	// for p := range md.tiles {
	// 	pos = append(pos, p)
	// }
	// return pos
}

func (md MoveDomain) HashBytes() []byte {
	return md.bitmap.HashBytes()
	// var res []byte
	// for _, pos := range SortedPositions(md.ListPosition()) {
	// 	res = append(res, byte(pos.X), byte(pos.Y))
	// }
	// return res
}

func (md MoveDomain) String() string {
	buf := &strings.Builder{}
	buf.WriteString("MoveDomain{")
	// for pos := range md.tiles {
	for _, pos := range md.bitmap.List() {
		buf.WriteString(pos.String())
	}
	buf.WriteString("}")
	return buf.String()
}

// func (md *MoveDomain) AddContacts(pos Pos) {
// 	md.contacts[pos] = struct{}{}
// }

// func (md MoveDomain) ContactsList() []Pos {
// 	var contacts []Pos
// 	for p := range md.contacts {
// 		contacts = append(contacts, p)
// 	}
// 	return contacts
// }
