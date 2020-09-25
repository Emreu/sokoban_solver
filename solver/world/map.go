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

var tileChars = map[byte]Tile{
	'#': TileWall,
	'@': TilePlayerStart,
	'+': TilePlayerStart,
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
				return Map{}, fmt.Errorf("error reading map @(%d,%d): %v", x, y, err)
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

		// TODO: validate line width
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
	// TODO: check if it's valid to have unmatching counts
	if boxCount != goalCount {
		err = fmt.Errorf("not equal numbers of boxes(%d) and goals(%d)", boxCount, goalCount)
	}
	if playerStartCount != 1 {
		err = fmt.Errorf("unexpected number of player start positions: %d", playerStartCount)
	}

	return Map{
		Width:  w,
		Height: len(tiles),
		StartX: pX,
		StartY: pY,
		Tiles:  tiles,
	}, err
}
