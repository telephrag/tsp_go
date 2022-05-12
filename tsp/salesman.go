package tsp

type Salesman struct {
	ID       uint64
	Copied   bool
	Path     []uint64 // list of node ids in order of travel
	Distance uint64   // length of entire path
	Count    uint64   // amount of nodes traveled (!) to
	Visited  uint64
}

func (sm *Salesman) Visit(nodeID uint64) {
	sm.Visited += (1 << nodeID)
}

func (sm *Salesman) HasVisited(nodeID uint64) bool {
	var bit uint64 = 1 << nodeID // get position of bit corresponding to given node
	mask := sm.Visited

	mask = mask & bit // XOR mask with bit
	mask = mask >> nodeID
	mask = mask << 63 >> 63

	return mask == 1
}

func (sm *Salesman) TailNode() uint64 {
	return sm.Path[sm.Count]
}
