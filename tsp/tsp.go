package tsp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"sync"
)

var graph []map[int]int // consider passing as parameters instead of storing globally

type Graph struct {
	Graph     []map[int]int `json:"adjList"`   // graph as adjancency list
	NodeCount int           `json:"nodeCount"` // amount of nodes
}

type Tsp struct {
	NodeCount int
	MinRoute  []int
	MinDist   int // current minimal path length
	Mu        sync.Mutex
}

func Init(inputPath string, maxOngoingTravels int, queueSize int) *Tsp {

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
	t.NodeCount = g.NodeCount // state
	t.MinRoute = make([]int, g.NodeCount)
	t.MinDist = math.MaxInt

	graph = g.Graph

	return &t
}
