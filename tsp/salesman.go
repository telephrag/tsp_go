package tsp

import (
	"fmt"
	"log"
	"sync/atomic"
)

type Salesman struct {
	ID       uint64
	Copied   uint32   // bool substiture
	Path     []uint64 // list of node ids in order of travel
	Distance uint64   // length of entire path
	Count    uint64   // amount of nodes traveled (!) to
	Visited  uint64   // bit mask representing visited nodes
}

func NewSalesman(t *Tsp) *Salesman {
	return &Salesman{
		ID:   atomic.AddUint64(&t.SalesmenCount, 1),
		Path: make([]uint64, t.NodeCount),
	}
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

func (sm *Salesman) Copy(t *Tsp) *Salesman {
	smc := &Salesman{}

	if sm.Copied == 1 {
		smc.ID = atomic.AddUint64(&t.SalesmenCount, 1)

	} else {
		atomic.StoreUint32(&(sm.Copied), 1)
		smc.ID = sm.ID // first copy of a salesman is considered the same salesman
	}

	smc.Path = make([]uint64, len(sm.Path))
	count := copy(smc.Path, sm.Path)
	if count != len(sm.Path) {
		fmt.Println(sm)
		log.Fatalln("Failed to copy path into salesman's clone")
	}

	smc.Copied = 0 // false
	smc.Distance = sm.Distance
	smc.Count = sm.Count
	smc.Visited = sm.Visited

	return smc
}
