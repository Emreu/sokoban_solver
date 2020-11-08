package world

// Find dead zones - positions from where it is impossible to move box to goal
func FindDeadZones(m Map) Bitmap {
	var deadCorners []Pos
	var deadZones Bitmap
	addDeadZone := func(x, y int) {
		pos := Pos{x, y}
		if !m.IsInside(pos) {
			return
		}
		switch m.AtPos(pos) {
		case TileWall, TileGoal, TilePlayerOnGoal, TileBoxOnGoal:
			// goal can't be deadzone
			return
		}
		// don't add if already added
		if deadZones.CheckBit(pos) {
			return
		}
		deadCorners = append(deadCorners, pos)
		deadZones.SetBit(pos)
	}
	// at first find dead corners
	for y, row := range m.Tiles {
		for x, tile := range row {
			if tile == TileWall {
				var up, down, left, right bool
				// check diagonal positions
				if m.At(x-1, y-1) == TileWall {
					up = true
					left = true
				}
				if m.At(x+1, y-1) == TileWall {
					up = true
					right = true
				}
				if m.At(x+1, y+1) == TileWall {
					right = true
					down = true
				}
				if m.At(x-1, y+1) == TileWall {
					left = true
					down = true
				}
				// add corresponding postition to dead zones if there is no wall or goal
				if up {
					addDeadZone(x, y-1)
				}
				if right {
					addDeadZone(x+1, y)
				}
				if down {
					addDeadZone(x, y+1)
				}
				if left {
					addDeadZone(x-1, y)
				}
			}
		}
	}
	// propagate dead corners
	for _, startPos := range deadCorners {
		for _, dir := range AllDirections {
			pos := startPos.MoveInDirection(dir)
			if m.AtPos(pos) == TileWall {
				continue
			}
			var leftOpen, rightOpen bool
			var deadCorridor []Pos
			// advance position & check left & right walls
		rayScan:
			for {
				// if out of map - stop
				if !m.IsInside(pos) {
					break rayScan
				}
				// if goal on path - it's not dead zone
				switch m.AtPos(pos) {
				case TileGoal, TileBoxOnGoal, TilePlayerOnGoal:
					break rayScan
				}
				// check left if it's not open yet
				if !leftOpen {
					leftPos := pos.MoveInDirection(dir.RotateCCW())
					if m.AtPos(leftPos) != TileWall {
						leftOpen = true
					}
				}
				// check right if it's not open yet
				if !rightOpen {
					rightPos := pos.MoveInDirection(dir.RotateCW())
					if m.AtPos(rightPos) != TileWall {
						rightOpen = true
					}
				}
				// if both open - stop propagation
				if leftOpen && rightOpen {
					break rayScan
				}
				// move ahead
				deadCorridor = append(deadCorridor, pos)
				pos = pos.MoveInDirection(dir)
				// if wall hit - stop propagation and add corridor to dead zone
				if m.AtPos(pos) == TileWall {
					for _, p := range deadCorridor {
						deadZones.SetBit(p)
					}
					break rayScan
				}
			}
		}
	}
	return deadZones
}
