package tsp

import (
	"sync"
)

func (t *Tsp) Solve() {
	sm := t.NewSalesman()
	sm.Path[0] = 0 // consider zero-node visited from the beginning
	sm.Visit(0)

	var wg sync.WaitGroup
	for nextNode := range graph[sm.TailNode()] {
		wg.Add(1)
		go func(smc *Salesman, node uint64) {
			defer wg.Done()
			t.travel(smc, node)
		}(t.CopySalesman(sm), nextNode)
	}
	wg.Wait()
}

func (t *Tsp) travel(sm *Salesman, node uint64) {
	sm.Distance += graph[sm.TailNode()][node]
	if sm.Distance > t.MinDist {
		return
	}

	sm.Visit(node)
	sm.Count++
	sm.Path[sm.Count] = node

	// if sm.ID == 1 {
	// 	fmt.Printf("t:  %2d %v %2d\n", sm.ID, sm.Path, sm.Distance)
	// }

	//fmt.Printf("t:   %6d %v %2d\n", sm.ID, sm.Path, sm.Distance)

	var wg sync.WaitGroup
	for nextNode := range graph[node] {
		wg.Add(1)
		go func(smc *Salesman, node uint64) {
			defer wg.Done()
			if !smc.HasVisited(node) {
				t.travel(smc, node)
			}
		}(t.CopySalesman(sm), nextNode)
	}
	wg.Wait()

	if sm.Count == t.NodeCount-1 {
		sm.Count++
		sm.Distance += graph[node][0]

		// if sm.Path[1] == 1 {
		// 	fmt.Printf("te: %6d %v %2d\n", sm.ID, sm.Path, sm.Distance)
		// }

		// fmt.Printf("te:  %6d %v %2d\n", sm.ID, sm.Path, sm.Distance)

		t.Mu.Lock()
		defer t.Mu.Unlock()
		if t.MinDist > sm.Distance {
			t.MinPath = sm.Path
			t.MinDist = sm.Distance
		}
	}
}

func (t *Tsp) travelEnd(sm *Salesman, node uint64) {
	sm.Distance += graph[sm.TailNode()][node]
	sm.Visit(node)
	sm.Count++
	sm.Path[sm.Count] = node

	// if sm.ID == 1 {
	// 	fmt.Printf("te:  %2d %v %2d\n", sm.ID, sm.Path, sm.Distance)
	// }
}
