package world

import (
	"bufio"
	"fmt"
	"io"
)

type Tile int

const (
	TileEmpty Tile = iota
	TileWall
	TilePlayerStart
	TilePlayerOnGoal
	TileBox
	TileBoxOnGoal
	TileGoal
)

type Map struct {
	Width  int
	Height int
	StartX int
	StartY int
	Tiles  [][]Tile
}

func (m Map) String() string {
	return fmt.Sprintf("Map %dx%d", m.Width, m.Height)
}

var tileChars = map[byte]Tile{
	'#': TileWall,
	'@': TilePlayerStart,
	'+': TilePlayerOnGoal,
	'$': TileBox,
	'*': TileBoxOnGoal,
	'.': TileGoal,
	' ': TileEmpty,
}

func tileFromChar(c byte) (Tile, error) {
	tile, ok := tileChars[c]
	if !ok {
		return Tile(-1), fmt.Errorf("unknown char: %c", c)
	}
	return tile, nil
}

func ReadMap(r io.Reader) (Map, error) {
	var w int
	var pX, pY int
	var tiles [][]Tile
	var boxCount, goalCount, playerStartCount int

	s := bufio.NewScanner(r)

	y := 0
	for s.Scan() {
		var line []Tile
		for x, c := range s.Bytes() {
			tile, err := tileFromChar(c)
			if err != nil {
				return Map{}, fmt.Errorf("%v @(%d,%d)", err, x, y)
			}
			line = append(line, tile)

			// calculate stats
			switch tile {
			case TileBox:
				boxCount++
			case TileBoxOnGoal:
				boxCount++
				goalCount++
			case TileGoal:
				goalCount++
			case TilePlayerStart:
				pX = x
				pY = y
				playerStartCount++
			case TilePlayerOnGoal:
				pX = x
				pY = y
				playerStartCount++
				goalCount++
			}
		}

		if len(line) == 0 {
			continue
		}

		if len(line) > w {
			w = len(line)
		}
		tiles = append(tiles, line)
		y++
	}

	var err error
	if boxCount != goalCount {
		err = fmt.Errorf("not equal numbers of boxes(%d) and goals(%d)", boxCount, goalCount)
	}
	if playerStartCount != 1 {
		err = fmt.Errorf("unexpected number of player start positions: %d", playerStartCount)
	}
	// fill incomplete rows with empty tile
	for y, row := range tiles {
		for len(row) < w {
			row = append(row, TileEmpty)
		}
		tiles[y] = row
	}

	return Map{
		Width:  w,
		Height: len(tiles),
		StartX: pX,
		StartY: pY,
		Tiles:  tiles,
	}, err
}

func (m Map) InitialBoxPositions() []Pos {
	var pos []Pos
	for y, row := range m.Tiles {
		for x, tile := range row {
			switch tile {
			case TileBox, TileBoxOnGoal:
				pos = append(pos, Pos{x, y})
			}
		}
	}
	return pos
}

func (m Map) At(x, y int) Tile {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return TileEmpty
	}
	return m.Tiles[y][x]
}

func (m Map) AtPos(p Pos) Tile {
	return m.At(p.X, p.Y)
}

func (m Map) IsInside(p Pos) bool {
	return p.X >= 0 && p.X < m.Width && p.Y >= 0 && p.Y < m.Height
}

func (m Map) StartPos() Pos {
	return Pos{X: m.StartX, Y: m.StartY}
}
