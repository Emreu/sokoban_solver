package world

import (
	"fmt"
	"log"
	"strings"
)

type cell struct {
	dist map[Pos]int
}

type MetricCalculator struct {
	cells [][]cell
}

func NewMetricCalculator(m Map) MetricCalculator {
	cells := make([][]cell, m.Height)
	for y := 0; y < m.Height; y++ {
		row := make([]cell, m.Width)
		for x := 0; x < m.Width; x++ {
			c := cell{
				dist: make(map[Pos]int),
			}
			switch m.At(x, y) {
			case TileGoal, TilePlayerOnGoal, TileBoxOnGoal:
				c.dist[Pos{x, y}] = 0
			}
			row[x] = c
		}
		cells[y] = row
	}
	mc := MetricCalculator{
		cells: cells,
	}
	mc.propagate(m)
	return mc
}

func (mc *MetricCalculator) propagate(m Map) {
	log.Print("Running metric propagation...")
	// run propagation until no new updates are made
	for {
		noUpdates := true
		for y, row := range mc.cells {
			for x, c := range row {
				curPos := Pos{x, y}
				if m.AtPos(curPos) == TileWall {
					continue
				}
				for _, dir := range []MoveDirection{MoveUp, MoveRight, MoveDown, MoveLeft} {
					p := curPos.MoveInDirection(dir)
					if !m.IsInside(p) {
						continue
					}
					// check if 2 step ahead is not wall, so it's possible to move box from there
					// if m.AtPos(p.MoveInDirection(dir)) == TileWall {
					// 	continue
					// }
					nc := mc.cells[p.Y][p.X]

					// loop over all goals in neighbour cell
					for pos, value := range nc.dist {
						if _, exists := c.dist[pos]; !exists {
							mc.cells[y][x].dist[pos] = value + 1
							noUpdates = false
						}
					}
				}
			}
		}
		if noUpdates {
			break
		}
	}

	log.Print(mc.String())
}

func (mc MetricCalculator) Evaluate(s State) int {
	// TODO: implement
	// for every box positions get corresponding cell
	// sort box on cells based on best possible value
	// eliminate goals double targeting by reassigning box to second best goals
	// sum total value
	return -1
}

func (mc MetricCalculator) String() string {
	buf := &strings.Builder{}
	buf.WriteString("Metric{\n")
	for y, row := range mc.cells {
		for x, c := range row {
			fmt.Fprintf(buf, "@(%d,%d)\n", x, y)
			for pos, value := range c.dist {
				fmt.Fprintf(buf, "->%s=%d\n", pos, value)
			}
		}
	}
	buf.WriteString("}")
	return buf.String()
}
