package tsp

import (
	"context"
	"time"
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
	}
	distance += graph[route[len(route)-1]][0] // return to zero-node

	t.Mu.Lock()
	if distance < t.MinDist {
		// fmt.Printf("Swap: %v <-> %v\n", t.MinRoute, route)
		t.MinDist = distance
		copy(t.MinRoute, route)
		copy(t.MinRoute, t.MinRoute[:1])
		// fmt.Printf("Swap: %v <-> %v\n", t.MinRoute, route)
	}
	t.Mu.Unlock()

}

func (t *Tsp) Solve() {
	if graph == nil {
		panic("graph is empty")
	}

	ctx, cancel := context.WithCancel(context.Background())

	jobChan := make(chan []int, config.RouteQueueSize) // limit amount of permutations at the same time
	workers := make(chan rune, config.MaxTravelers)    // limit amount of workers

	nodeSet := make([]int, t.NodeCount)
	for i := 1; i < t.NodeCount; i++ {
		nodeSet[i] = i
	}

	go func() {
		for r := make([]int, len(nodeSet)); r[0] < len(r); nextPermutation(r) {
			jobChan <- getPermutation(nodeSet, r)
		}
		// create permutation in a seperate goroutine
		cancel() // cancel the context once done
	}()

	for {
		select {
		case routeToCalc := <-jobChan:
			workers <- 1 // book the slot once received a job
			go func() {
				t.calcDistance(routeToCalc)
				<-workers // free the slot when done
			}()
		case <-ctx.Done(): // wont be checked if previouse case is proced
			for len(workers) > 0 {
				time.Sleep(time.Millisecond * 10) // wait for the last workers to finish
			}
			return // e.g. when there is still job in jobChan
		}
	}
}
