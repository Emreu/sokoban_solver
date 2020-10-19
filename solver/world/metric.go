package world

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

type cell struct {
	dist map[Pos]int
}

type MetricCalculator struct {
	cells [][]cell
}

func NewMetricCalculator(m Map, deadzones MoveDomain) MetricCalculator {
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
	mc.propagate(m, deadzones)
	return mc
}

func (mc *MetricCalculator) propagate(m Map, deadzones MoveDomain) {
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
				if deadzones.HasPosition(curPos) {
					continue
				}
				for _, dir := range []MoveDirection{MoveUp, MoveRight, MoveDown, MoveLeft} {
					p := curPos.MoveInDirection(dir)
					if !m.IsInside(p) {
						continue
					}
					pushPos := curPos.MoveAgainstDirection(dir)
					// check if there is no wall at push position, so it's possible to move box
					if !m.IsInside(pushPos) {
						continue
					}
					if m.AtPos(pushPos) == TileWall {
						continue
					}
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

	// log.Print(mc.String())
}

type boxGoal struct {
	min     int
	minGoal Pos
	metrics map[Pos]int
}

func (bg *boxGoal) ExcludeGoal(g Pos) bool {
	delete(bg.metrics, g)
	// if no more goals available - state is failed
	if len(bg.metrics) == 0 {
		return false
	}
	bg.min = 99999
	for pos, m := range bg.metrics {
		if m < bg.min {
			bg.min = m
			bg.minGoal = pos
		}
	}
	return true
}

func (mc MetricCalculator) Evaluate(s State) (int, error) {
	var sortList []boxGoal
	var total = 0

boxes:
	for _, pos := range s.BoxPositions {
		bg := boxGoal{
			min:     99999,
			metrics: make(map[Pos]int),
		}
		for goalPos, metric := range mc.cells[pos.Y][pos.X].dist {
			if metric == 0 {
				// skip this box from sorting
				continue boxes
			}
			bg.metrics[goalPos] = metric
			if metric < bg.min {
				bg.min = metric
				bg.minGoal = goalPos
			}
		}
		sortList = append(sortList, bg)
	}

	// eliminate boxes from sort list
	for len(sortList) > 0 {
		// sort boxes
		sort.SliceStable(sortList, func(i, j int) bool {
			return sortList[i].min < sortList[j].min
		})

		// remove first box and add up metric and disable goal
		best := sortList[0]
		total += best.min
		sortList = sortList[1:len(sortList)]
		for i := range sortList {
			ok := sortList[i].ExcludeGoal(best.minGoal)
			if !ok {
				return -1, fmt.Errorf("failed state")
			}
		}
	}

	return total, nil
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
