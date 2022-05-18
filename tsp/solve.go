package tsp

import (
	"sync"
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

func (t *Tsp) calcDistance(route []int, startNode int) {

	distance := graph[0][startNode]        // add distance from zero-node to the first node in route
	distance += graph[startNode][route[0]] // account for startNode
	for i := 0; i < len(route)-1; i++ {    // calculate weight of the rest of the permutation
		distance += graph[route[i]][route[i+1]]
		if distance > t.MinDist {
			return
		}
	}
	distance += graph[route[len(route)-1]][0] // return to zero-node

	t.Mu.Lock()
	if distance < t.MinDist {
		t.MinDist = distance
		copy(t.MinRoute[2:], route) // copy route into t.MinRoute with shift by two
		t.MinRoute[1] = startNode   // the zero-elem will be 0 and the first will be startNode
	}
	t.Mu.Unlock()
}

func (t *Tsp) travel(startNode int) {
	nodeSet := make([]int, t.NodeCount-2)
	j := 0
	for i := 1; i < t.NodeCount; i++ {
		if i != startNode {
			nodeSet[j] = i
			j++
		}
	}

	for r := make([]int, len(nodeSet)); r[0] < len(r); nextPermutation(r) {
		t.calcDistance(
			getPermutation(nodeSet, r),
			startNode,
		)
	}
}

func (t *Tsp) Solve() {
	if graph == nil {
		panic("graph is empty")
	}

	nodeSet := make([]int, t.NodeCount-1)
	for i := 0; i < len(nodeSet); i++ {
		nodeSet[i] = i + 1
	}

	var wg sync.WaitGroup
	wg.Add(len(nodeSet))
	for _, startNode := range nodeSet {
		go func(sn int) {
			t.travel(sn)
			wg.Done()
		}(startNode)
	}
	wg.Wait()
}
