package world

import "hash/maphash"

type State struct {
	MoveDomain
	BoxPositions []Pos
}

func (s State) Hash() uint64 {
	var h maphash.Hash
	// hash box positions
	for _, pos := range SortedPositions(s.BoxPositions) {
		h.WriteByte(byte(pos.X))
		h.WriteByte(byte(pos.Y))
	}
	// hash move domain
	for _, pos := range SortedPositions(s.MoveDomain.ListPosition()) {
		h.WriteByte(byte(pos.X))
		h.WriteByte(byte(pos.Y))
	}
	return h.Sum64()
}
