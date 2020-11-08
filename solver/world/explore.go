package world

import "log"

func Explore(m Map, boxes PosList, start Pos, deadZones, hLock, vLock Bitmap) (Bitmap, []BoxMove) {
	var domain Bitmap
	var boxBitmap Bitmap
	var moves []BoxMove

	for _, pos := range boxes {
		boxBitmap.SetBit(pos)
	}

	boxIndex := func(p Pos) int {
		for i, pos := range boxes {
			if p == pos {
				return i
			}
		}
		return -1
	}

	var posQueue = []Pos{start}
	var comeDirections = []MoveDirection{-1}
	var distances = []int{0}

queue:
	for i := 0; i < len(posQueue); i++ {
		p := posQueue[i]
		d := comeDirections[i]
		dist := distances[i]
		if !m.IsInside(p) {
			continue
		}
		if m.AtPos(p) == TileWall {
			continue
		}
		backward := d.Opposite()
		// if we hit a box - check if move is possible and save
		if boxBitmap.CheckBit(p) {
			dstPos := p.MoveInDirection(d)
			// check if destination position is available
			if m.AtPos(dstPos) == TileWall {
				continue
			}
			if boxBitmap.CheckBit(dstPos) {
				continue
			}
			// check destination isn't in dead zone
			if deadZones.CheckBit(dstPos) {
				continue
			}
			// for quick deadlock detection check if there is neighbour box with same h/v lock as target position
			// and destination pos isn't goal - move to goal is ok even if deadlock
			dstHLock := hLock.CheckBit(dstPos)
			dstVLock := vLock.CheckBit(dstPos)
			dstGoal := m.AtPos(dstPos) == TileGoal || m.AtPos(dstPos) == TileBoxOnGoal || m.AtPos(dstPos) == TilePlayerOnGoal
			if (dstHLock || dstVLock) && !dstGoal {
				// check only if next step to move is blocked by box or wall
				aheadPos := dstPos.MoveInDirection(d)
				if m.AtPos(aheadPos) == TileWall || boxBitmap.CheckBit(aheadPos) {
					for _, dir := range AllDirections {
						if dir == backward {
							continue
						}
						n := dstPos.MoveInDirection(dir)
						if !boxBitmap.CheckBit(n) {
							continue
						}
						if dstHLock && hLock.CheckBit(n) {
							log.Printf("Move discarded due to horizontal lock")
							continue queue
						}
						if dstVLock && vLock.CheckBit(n) {
							log.Printf("Move discarded due to vertical lock")
							continue queue
						}
					}
				}
			}
			// save move
			move := BoxMove{
				Direction: d,
				Distance:  dist,
				BoxIndex:  boxIndex(p),
			}
			moves = append(moves, move)
			continue
		}
		if domain.CheckBit(p) {
			continue
		}
		domain.SetBit(p)
		for _, dir := range AllDirections {
			if dir == backward {
				continue
			}
			n := p.MoveInDirection(dir)
			posQueue = append(posQueue, n)
			comeDirections = append(comeDirections, dir)
			distances = append(distances, dist+1)
		}
	}

	return domain, moves
}

func ExploreDomainOnly(m Map, boxes PosList, start Pos) Bitmap {
	var domain Bitmap
	var boxBitmap Bitmap

	for _, pos := range boxes {
		boxBitmap.SetBit(pos)
	}

	var posQueue = []Pos{start}
	var comeDirections = []MoveDirection{-1}

	for i := 0; i < len(posQueue); i++ {
		p := posQueue[i]
		if !m.IsInside(p) {
			continue
		}
		if m.AtPos(p) == TileWall {
			continue
		}
		if boxBitmap.CheckBit(p) {
			continue
		}
		if domain.CheckBit(p) {
			continue
		}
		domain.SetBit(p)
		backward := comeDirections[i].Opposite()
		for _, dir := range AllDirections {
			if dir == backward {
				continue
			}
			n := p.MoveInDirection(dir)
			posQueue = append(posQueue, n)
			comeDirections = append(comeDirections, dir)
		}
	}

	return domain
}
