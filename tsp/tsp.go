package tsp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"sync"
	"sync/atomic"
)

var graph []map[uint64]uint64 // consider passing as parameters instead of storing globally

type Graph struct {
	Graph     []map[uint64]uint64 `json:"adjList"`   // graph as adjancency list
	NodeCount uint64              `json:"nodeCount"` // amount of nodes
}

type Tsp struct {
	NodeCount     uint64
	MinPath       []uint64
	MinDist       uint64 // current minimal path length
	SalesmenCount uint64
	Wg            sync.WaitGroup
	Mu            sync.Mutex
}

func Init(inputPath string) *Tsp {

	content, err := ioutil.ReadFile(inputPath)
	if err != nil {
		log.Fatalln(err)
	}

	g := Graph{}
	err = json.Unmarshal(content, &g)
	if err != nil {
		log.Fatalln(err)
	}

	t := Tsp{}
	t.NodeCount = g.NodeCount
	t.MinPath = make([]uint64, g.NodeCount)
	t.MinDist = math.MaxUint64

	graph = g.Graph

	return &t
}

func (t *Tsp) NewSalesman() *Salesman {
	return &Salesman{
		ID:   atomic.AddUint64(&t.SalesmenCount, 1),
		Path: make([]uint64, t.NodeCount),
	}
}

func (t *Tsp) CopySalesman(sm *Salesman) *Salesman {
	smc := &Salesman{}

	if sm.Copied {
		smc.ID = atomic.AddUint64(&t.SalesmenCount, 1)

		smc.Path = make([]uint64, len(sm.Path))
		count := copy(smc.Path, sm.Path)
		if count != len(sm.Path) {
			fmt.Println(sm)
			log.Fatalln("Failed to copy path into salesman's clone")
		}

	} else {
		smc.ID = sm.ID // first copy of a salesman is considered the same salesman

		smc.Path = sm.Path
		sm.Copied = true // now salesman is considered copied
	}

	smc.Copied = false
	smc.Distance = sm.Distance
	smc.Count = sm.Count
	smc.Visited = sm.Visited

	return smc
}

func (t *Tsp) RollbackCopy(sm *Salesman, returnBy uint64) *Salesman {

	smc := t.NewSalesman()
	count := copy(smc.Path, sm.Path[:returnBy])
	if count != len(sm.Path[:returnBy]) {
		log.Fatalln("Failed to copy rolled back path")
	}

	return smc
}
