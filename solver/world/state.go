package world

import "hash/maphash"

type State struct {
	MoveDomain
	BoxPositions []Pos
}

var hash maphash.Hash

func (s State) Hash() uint64 {
	hash.Reset()
	// hash box positions
	for _, pos := range SortedPositions(s.BoxPositions) {
		hash.WriteByte(byte(pos.X))
		hash.WriteByte(byte(pos.Y))
	}
	// hash move domain
	hash.Write(s.MoveDomain.HashBytes())
	return hash.Sum64()
}
