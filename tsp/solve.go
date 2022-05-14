package tsp

func (t *Tsp) Solve() {
	sm := NewSalesman(t)
	sm.Visit(0) // sets bit corresponding to given node in bitmask

	for nextNode := range graph[0] {
		t.RouteQueue <- true // book slot in root queue
		go t.travel(sm.Copy(t), nextNode)
	}
}

func (t *Tsp) travel(sm *Salesman, node uint64) {
	sm.Distance += graph[sm.TailNode()][node] // increase traveled distance
	if sm.Distance > t.MinDist {              // terminate if traveled distance is bigger than current minimal
		<-t.RouteQueue
		return
	}

	sm.Visit(node)
	sm.Count++               // increase amount of nodes traveled
	sm.Path[sm.Count] = node // add node to path

	if sm.Count == t.NodeCount-1 { // stop if t.NodeCount - 1 nodes traveled
		sm.Count++
		sm.Distance += graph[node][0] // return to zero-node

		t.Mu.Lock()
		if t.MinDist > sm.Distance { // set new min distance and path if they are shorter
			t.MinPath = sm.Path
			t.MinDist = sm.Distance
		}
		t.Mu.Unlock()
	}

	<-t.RouteQueue // free the slot for routes

	for nextNode := range graph[node] {
		if !sm.HasVisited(node) {
			t.RouteQueue <- true
			go t.travel(sm.Copy(t), nextNode)
		}
	}
}
