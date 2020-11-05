package world

type Frontier struct{}

func (f *Frontier) Add(n *Node) {

}

func (f *Frontier) Pick() (*Node, error) {
	return nil, nil
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
