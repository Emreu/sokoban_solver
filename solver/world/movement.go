package world

type MoveDirection int

const (
	MoveUp MoveDirection = iota
	MoveRight
	MoveDown
	MoveLeft
)

type MoveDomain struct {
	tiles map[Pos]struct{} // simple implementation
	// contacts map[Pos]struct{}
}

func NewMoveDomain() MoveDomain {
	return MoveDomain{
		tiles: make(map[Pos]struct{}),
		// contacts: make(map[Pos]struct{}),
	}
}

// NewMoveDomainFromMap generate move domain from Map (considering only walls), box positions and starting position
func NewMoveDomainFromMap(m Map, boxPositions []Pos, start Pos) MoveDomain {
	var perimeter = make(map[Pos]struct{})
	perimeter[start] = struct{}{}
	domain := NewMoveDomain()

	var boxPos = make(map[Pos]struct{})
	for _, pos := range boxPositions {
		boxPos[pos] = struct{}{}
	}

	var nextPerimeter = make(map[Pos]struct{})
	for {
		for pos := range perimeter {
			// skip if this tile is wall
			if m.AtPos(pos) == TileWall {
				continue
			}
			// skip if tile is occupied by box
			if _, occupied := boxPos[pos]; occupied {
				// but save contact position
				// domain.AddContacts(pos)
				continue
			}
			// add tile to domain
			domain.AddPosition(pos)
			// schedule neighbour tiles for next perimeter
			for _, p := range pos.Neighbours() {
				// skip if out of map
				if !m.IsInside(p) {
					continue
				}
				// skip if already processing
				if _, processing := perimeter[p]; processing {
					continue
				}
				// skip if already scheduled
				if _, scheduled := nextPerimeter[p]; scheduled {
					continue
				}
				// skip if already in domain
				if domain.HasPosition(p) {
					continue
				}
				// finally add to next perimeter
				nextPerimeter[p] = struct{}{}
			}
		}
		if len(nextPerimeter) == 0 {
			break
		}
		perimeter = nextPerimeter
		nextPerimeter = make(map[Pos]struct{})
	}

	return domain
}

func (md *MoveDomain) HasPosition(pos Pos) bool {
	_, exists := md.tiles[pos]
	return exists
}

func (md *MoveDomain) AddPosition(pos Pos) {
	md.tiles[pos] = struct{}{}
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
