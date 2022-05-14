package main

import (
	"fmt"
	"time"
	"tsp/config"
	"tsp/tsp"
)

func measureExecTime(f func()) {
	start := time.Now()
	f()
	elapsed := time.Since(start)
	fmt.Printf("\nfinished in %d us\n", elapsed.Microseconds())
}

func main() {
	t := tsp.Init("input_a.json", config.MaxRoutes)
	t.Solve()

	fmt.Printf("Solution: %v %d\n", t.MinPath, t.MinDist)
}
