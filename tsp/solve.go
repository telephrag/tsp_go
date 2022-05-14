package tsp

func (t *Tsp) Solve() {
	sm := NewSalesman(t)
	sm.Path[0] = 0 // consider zero-node visited from the beginning
	sm.Visit(0)

	for nextNode := range graph[sm.TailNode()] {
		t.RouteQueue <- true
		go t.travel(sm.Copy(t), nextNode)
	}
}

func (t *Tsp) travel(sm *Salesman, node uint64) {
	sm.Distance += graph[sm.TailNode()][node]
	if sm.Distance > t.MinDist {
		<-t.RouteQueue
		return
	}

	sm.Visit(node)
	sm.Count++
	sm.Path[sm.Count] = node

	if sm.Count == t.NodeCount-1 {
		sm.Count++
		sm.Distance += graph[node][0]

		t.Mu.Lock()
		if t.MinDist > sm.Distance {
			t.MinPath = sm.Path
			t.MinDist = sm.Distance
		}
		t.Mu.Unlock()
	}

	<-t.RouteQueue
	for nextNode := range graph[node] {
		if !sm.HasVisited(node) {
			t.RouteQueue <- true
			go t.travel(sm.Copy(t), nextNode)
		}
	}
}
