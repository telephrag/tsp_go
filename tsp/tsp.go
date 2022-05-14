package tsp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"sync"
)

var graph []map[uint64]uint64

type Graph struct {
	Graph     []map[uint64]uint64 `json:"adjList"`   // graph as adjancency list
	NodeCount uint64              `json:"nodeCount"` // amount of nodes
}

type Tsp struct {
	NodeCount     uint64    // amount of nodes in travelled graph
	MinPath       []uint64  // current shortest route
	MinDist       uint64    // current shortest route length
	SalesmenCount uint64    // used for issuing IDs to salesmen
	RouteQueue    chan bool // used to limit amount of concurrent travels
	Mu            sync.Mutex
}

func Init(inputPath string, maxRoutes int) *Tsp {

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
	t.RouteQueue = make(chan bool, maxRoutes)

	graph = g.Graph

	return &t
}
