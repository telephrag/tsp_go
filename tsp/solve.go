package tsp

import (
	"context"
	"sync"
	"tsp/config"
)

func nextPermutation(route []int) {
	for i := len(route) - 1; i >= 0; i-- {
		if i == 0 || route[i] < len(route)-i-1 {
			route[i]++
			return
		}
		route[i] = 0
	}
}

func getPermutation(original, route []int) []int {
	result := append([]int{}, original...)
	for i, v := range route {
		result[i], result[i+v] = result[i+v], result[i]
	}
	return result
}

func (t *Tsp) calcDistance(route []int) {

	distance := graph[0][route[0]] // add distance from zero-node to the first node in route
	for i := 0; i < len(route)-1; i++ {
		distance += graph[route[i]][route[i+1]]
		// if distance > t.MinDist {
		// 	return
		// }
	}
	distance += graph[route[len(route)-1]][0] // return to zero-node

	if distance < t.MinDist {
		t.Mu.Lock()
		// fmt.Printf("Swap: %v <-> %v\n", t.MinRoute, route)
		t.MinDist = distance
		copy(t.MinRoute, route)
		copy(t.MinRoute, t.MinRoute[:1])
		// fmt.Printf("Swap: %v <-> %v\n", t.MinRoute, route)
		t.Mu.Unlock()
	}

}

func (t *Tsp) travel(routeChan chan []int, ctx context.Context) {
	for {
		select {
		case routeToCalc := <-routeChan: // Possible datarace hence decreased performance
			t.calcDistance(routeToCalc)
			continue
		default:
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (t *Tsp) Solve() {
	if graph == nil {
		panic("graph is empty")
	}

	ctx, cancel := context.WithCancel(context.Background())
	routeChan := make(chan []int, config.RouteQueueSize) // limit amount of permutations at the same time

	nodeSet := make([]int, t.NodeCount)
	for i := 1; i < t.NodeCount; i++ {
		nodeSet[i] = i
	}

	// consider turning into buffer
	// consider it to evenly spread routes between travellers
	go func() {
		for r := make([]int, len(nodeSet)); r[0] < len(r); nextPermutation(r) {
			routeChan <- getPermutation(nodeSet, r)
		}
		// create permutation in a seperate goroutine
		cancel() // cancel the context once done
	}()

	var wg sync.WaitGroup
	wg.Add(config.MaxTravelers)
	for i := 0; i < config.MaxTravelers; i++ {
		go func() {
			t.travel(routeChan, ctx)
			wg.Done()
		}()
	}
	wg.Wait()
}
