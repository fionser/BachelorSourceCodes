package main

import (
	"FP_Logic"
	"log"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	aggregator := FP_Logic.NewAggregator()
	aggregator.Aggregate()
	log.Println(aggregator.AverageRadius())
}
