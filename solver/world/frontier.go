package world

import (
	"log"
	"sort"
)

const batchSize = 10

type Frontier struct {
	nodes []*Node
}

func (f *Frontier) Add(n *Node) {
	f.nodes = append(f.nodes, n)
}

func (f Frontier) HasNodes() bool {
	return len(f.nodes) > 0
}

func (f *Frontier) GetBatch() []*Node {
	// sort
	sort.Slice(f.nodes, func(i, j int) bool {
		return f.nodes[i].Metric < f.nodes[j].Metric
	})
	index := batchSize
	if index > len(f.nodes) {
		index = len(f.nodes)
	}
	top := f.nodes[:index]
	remaining := f.nodes[index:]
	log.Printf("Fetching from frontier %d nodes (%d remaining), best metric: %d, worst: %d", len(top), len(remaining), top[0].Metric, top[len(top)-1].Metric)

	f.nodes = remaining
	return top
}

func mergeSorted(base, addition []*Node) []*Node {
	var result = make([]*Node, len(base)+len(addition))
	var i, j, k = 0, 0, 0
	for i < len(base) && j < len(addition) {
		if base[i].Metric < addition[j].Metric {
			result[k] = base[i]
			i++
		} else {
			result[k] = addition[j]
			j++
		}
		k++
	}
	if i < len(base) {
		copy(result[k:], base[i:])
	} else {
		copy(result[k:], addition[j:])
	}
	return result
}

/*
for len(exploreFrontier) != 0 {
		if err := c.Err(); err != nil {
			return fmt.Errorf("timed out or canceled")
		}
		var nextFrontier []*Node
		log.Printf("Exploring frontier #%d of %d nodes (%d postponed)...", step, len(exploreFrontier), len(postponedFrontier))
		for _, node := range exploreFrontier {
			newNodes := s.exploreNode(node)
			for _, n := range newNodes {
				// check if solution was found
				if s.isSolution(n.Boxes) {
					log.Print("Solution found!")
					s.solution = n
					return nil
				}
				nextFrontier = append(nextFrontier, n)
			}
		}

		// make next frontier
		fullFrontier := mergeSorted(postponedFrontier, nextFrontier)
		if len(fullFrontier) > ExploreBatch {
			exploreFrontier = fullFrontier[:ExploreBatch]
			postponedFrontier = fullFrontier[ExploreBatch:]
		} else {
			exploreFrontier = fullFrontier
			postponedFrontier = []*Node{}
		}

		if len(exploreFrontier) > 0 {
			log.Printf("Best metric: %d", exploreFrontier[0].Metric)
		}

		step++
	}
*/
