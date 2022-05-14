package tsp

import "fmt"

func (t *Tsp) Solve() {
	if graph == nil {
		panic("graph is empty")
	}

	// creating initial route that is basically an ascending listing of 0..n-2 nodes
	sm := t.NewSalesman()
	sm.Visit(0)
	// at the end of execution sm.Count (equals to len(graph)) is number of nodes,
	// reducing it by mw ill result in m nodes from the right becoming null
	for i := uint64(1); i < uint64(len(graph))-2; i++ {
		sm.Distance += graph[sm.Path[sm.Count]][i] // add distance from last travelled node to one with nodeId
		sm.Count++
		sm.Path[i] = i
		sm.Visit(i)
		fmt.Printf("%b\n", sm.Visited)
	}

	fmt.Printf("Initial route: %v %d %d\n", sm.Path, sm.Distance, sm.Count)

	// fmt.Println("Testing UnVisit()")
	// for i := range sm.Path {
	// 	sm.UnVisit(uint64(i))
	// 	fmt.Printf("%b\n", sm.Visited)
	// }

	// smc := t.RollbackCopy(sm, sm.Count-1)
	// fmt.Printf("Testing RollbackCopy: %v\n", smc.Path)

	threadCount := 6
	for offset := 3; offset != len(graph)-2; offset++ { // would overflow occur if we do 0 - 1 on uint64?
		fmt.Printf("%3d %3d\n", offset, threadCount)
		threadCount *= offset
	}

}
