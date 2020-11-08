package world

func CheckDeadlock(m Map, boxes PosList, move BoxMove) bool {

	return false
}

func FindMoveLocks(m Map, skip Bitmap) (Bitmap, Bitmap) {
	var vLock, hLock Bitmap
	addLock := func(x, y int, h bool) {
		pos := Pos{x, y}
		if !m.IsInside(pos) {
			return
		}
		if skip.CheckBit(pos) {
			return
		}
		if m.AtPos(pos) == TileWall {
			return
		}
		if h {
			hLock.SetBit(pos)
		} else {
			vLock.SetBit(pos)
		}
	}
	for y, row := range m.Tiles {
		for x, tile := range row {
			if tile == TileWall {
				addLock(x-1, y, true)
				addLock(x+1, y, true)
				addLock(x, y-1, false)
				addLock(x, y+1, false)
			}
		}
	}
	return hLock, vLock
}
