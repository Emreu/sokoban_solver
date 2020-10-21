package world

import (
	"fmt"
	"sort"
)

func manhattan(from, to Pos) int {
	dist := 0
	if from.X > to.X {
		dist += from.X - to.X
	} else {
		dist += to.X - from.X
	}
	if from.Y > to.Y {
		dist += from.Y - to.Y
	} else {
		dist += to.Y - from.Y
	}
	return dist
}

func FindDirections(domain MoveDomain, from, to Pos) ([]MoveDirection, error) {
	if from == to {
		return []MoveDirection{}, nil
	}
	// check trivial one step move
	if manhattan(from, to) == 1 {
		for _, dir := range AllDirections {
			if from.MoveInDirection(dir) == to {
				return []MoveDirection{dir}, nil
			}
		}
	}
	// do normal search
	explore := []Pos{from}
	comeFrom := make(map[Pos]MoveDirection)

	traceback := func(final Pos) []MoveDirection {
		var path []MoveDirection
		pos := final
		for {
			dir, exists := comeFrom[pos]
			if !exists {
				break
			}
			pos = pos.MoveAgainstDirection(dir)
			path = append([]MoveDirection{dir}, path...)
		}
		return path
	}

	g := make(map[Pos]int)
	f := make(map[Pos]int)

	// addToExplore := func(p Pos) {
	// 	if len(explore) == 0 {
	// 		explore = []Pos{p}
	// 		return
	// 	}
	// 	// do binary search to find insertion index
	// 	threshold := f[p]
	// 	i := 0
	// 	j := len(explore) - 1
	// 	for i < j {
	// 		c := (i + j) / 2
	// 		if f[explore[c]] < threshold {
	// 			i = c
	// 		} else {
	// 			j = c
	// 		}
	// 	}
	// 	explore = append(explore, Pos{})
	// 	copy(explore[i+1:], explore[i:])
	// 	explore[i] = p
	// }

	g[from] = 0
	f[from] = manhattan(from, to)

	for len(explore) > 0 {
		// get best position from explore
		current := explore[0]
		// if it is destination - return path
		if current == to {
			return traceback(current), nil
		}
		explore = explore[1:]

		// probe all directions
		for _, dir := range AllDirections {
			// check if we stay inside movement domain
			neighbour := current.MoveInDirection(dir)
			if !domain.HasPosition(neighbour) {
				continue
			}
			// check if better path exists
			total := g[current] + 1
			if oldtotal, exists := g[neighbour]; exists && oldtotal < total {
				continue
			}
			// store this step
			comeFrom[neighbour] = dir
			g[neighbour] = total
			f[neighbour] = total + manhattan(neighbour, to)
			// add to explore
			// addToExplore(neighbour)
			// XXX: temporary fix until bsearch repair
			explore = append(explore, neighbour)
		}

		// TODO: remove after bsearch repair
		sort.Slice(explore, func(i, j int) bool {
			return f[explore[i]] < f[explore[j]]
		})
	}

	return nil, fmt.Errorf("path not exists")
}
