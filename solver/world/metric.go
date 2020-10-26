package world

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

type MetricCalculator struct {
	field      [][][]int // field [Y][X][goal_index] -> metric
	min        [][]int   // minimal metric value for position
	minIndex   [][]int   // minimal metric goal index for position
	goalsCount int
}

func initialMetrics(goals []Pos, p Pos) []int {
	res := make([]int, len(goals))
	for i := range goals {
		if goals[i] == p {
			res[i] = 0
		} else {
			res[i] = -1
		}
	}
	return res
}

func NewMetricCalculator(m Map, deadzones Bitmap) MetricCalculator {
	goals := m.GoalsPositions()
	field := make([][][]int, m.Height)
	for y := 0; y < m.Height; y++ {
		row := make([][]int, m.Width)
		for x := 0; x < m.Width; x++ {
			row[x] = initialMetrics(goals, Pos{x, y})
		}
		field[y] = row
	}
	mc := MetricCalculator{
		field:      field,
		goalsCount: len(goals),
	}
	mc.propagate(m, deadzones)
	mc.initMin()
	return mc
}

func (mc *MetricCalculator) propagate(m Map, deadzones Bitmap) {
	log.Print("Running metric propagation...")
	// run propagation until no new updates are made
	for {
		noUpdates := true
		for y, row := range mc.field {
			for x, metrics := range row {
				curPos := Pos{x, y}
				if m.AtPos(curPos) == TileWall {
					continue
				}
				if deadzones.CheckBit(curPos) {
					continue
				}
				for _, dir := range AllDirections {
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
					nMetrics := mc.field[p.Y][p.X]

					// loop over all goals in neighbour metrics
					for i := range nMetrics {
						if metrics[i] == -1 && nMetrics[i] >= 0 {
							metrics[i] = nMetrics[i] + 1
							noUpdates = false
						}
					}
				}
				mc.field[y][x] = metrics
			}
		}
		if noUpdates {
			break
		}
	}

	// log.Print(mc.String())
}

func (mc *MetricCalculator) initMin() {
	log.Print("Finding metrics min values...")
	mc.min = make([][]int, len(mc.field))
	mc.minIndex = make([][]int, len(mc.field))
	for y, row := range mc.field {
		minRow := make([]int, len(row))
		minIndexRow := make([]int, len(row))
		for x, metrics := range row {
			minIndex := -1
			for i, value := range metrics {
				if value == -1 {
					continue
				}
				if minIndex == -1 {
					minIndex = i
				} else if value < metrics[minIndex] {
					minIndex = i
				}
			}
			minIndexRow[x] = minIndex
			if minIndex >= 0 {
				minRow[x] = metrics[minIndex]
			}
		}
		mc.min[y] = minRow
		mc.minIndex[y] = minIndexRow
	}
}

type boxGoal struct {
	metrics  []int
	min      int
	minIndex int
}

func (mc MetricCalculator) Evaluate(s State) (int, error) {
	var sortList []boxGoal
	var total = 0
	boundGoals := make([]bool, mc.goalsCount)

	for _, pos := range s.BoxPositions {
		metrics := mc.field[pos.Y][pos.X]
		min := mc.min[pos.Y][pos.X]
		minIndex := mc.minIndex[pos.Y][pos.X]

		if minIndex == -1 {
			return -1, fmt.Errorf("unresolvable position")
		}

		if min == 0 {
			// skip this box from sorting
			boundGoals[minIndex] = true
			continue
		}

		sortList = append(sortList, boxGoal{
			metrics:  metrics,
			min:      min,
			minIndex: minIndex,
		})
	}

	// eliminate boxes from sort list
	for len(sortList) > 0 {
		// update min
		for i, item := range sortList {
			// do if minimal goal is bound
			if boundGoals[item.minIndex] {
				minIndex := -1
				for j, value := range item.metrics {
					// skip already bound goals
					if boundGoals[j] {
						continue
					}
					if value == -1 {
						continue
					}
					if minIndex == -1 {
						minIndex = j
					} else if value < item.metrics[minIndex] {
						minIndex = j
					}
				}
				if minIndex == -1 {
					return -1, fmt.Errorf("unresolvable position")
				}
				sortList[i].minIndex = minIndex
				sortList[i].min = item.metrics[minIndex]
			}
		}

		// sort boxes
		sort.SliceStable(sortList, func(i, j int) bool {
			return sortList[i].min < sortList[j].min
		})

		// remove first box and add up metric and disable goal
		best := sortList[0]
		total += best.min
		sortList = sortList[1:len(sortList)]
		boundGoals[best.minIndex] = true
	}

	return total, nil
}

func (mc MetricCalculator) String() string {
	buf := &strings.Builder{}
	buf.WriteString("Metric{\n")
	for y, row := range mc.field {
		for x, metrics := range row {
			fmt.Fprintf(buf, "@(%d,%d)\n", x, y)
			for i, value := range metrics {
				if value == -1 {
					continue
				}
				fmt.Fprintf(buf, "->%d=%d\n", i, value)
			}
		}
	}
	buf.WriteString("}")
	return buf.String()
}
