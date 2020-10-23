package world

import (
	"hash/fnv"
)

type State struct {
	MoveDomain
	BoxPositions []Pos
}

var hash = fnv.New64()

func (s State) Hash() uint64 {
	hash.Reset()
	// hash box positions
	for _, pos := range SortedPositions(s.BoxPositions) {
		hash.Write([]byte{byte(pos.X)})
		hash.Write([]byte{byte(pos.Y)})
	}
	// hash move domain
	hash.Write(s.MoveDomain.HashBytes())
	return hash.Sum64()
}
