package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

func addHarmonic(totalSoFar float64, i int) float64 {
	return totalSoFar + 1.0/float64(i)
}

func harmonicRange(start, end int) float64 {
	sum := 0.0
	for i := start; i <= end; i++ {
		sum = addHarmonic(sum, i)
	}
	return sum
}

var partialSums []float64
var wg sync.WaitGroup

func calcPartialSum(goroutine, start, end int) {
	defer wg.Done()
	sum := harmonicRange(start, end)
	partialSums[goroutine] = sum
	// fmt.Printf("goroutine %v for n=%v..%v got %v\n", goroutine, start, end, sum)
}

func main() {
	n, _ := strconv.Atoi(os.Args[1])
	goroutines, err := strconv.Atoi(os.Args[2])
	if err != nil {
		goroutines = runtime.NumCPU()
	}
	fmt.Printf("\n\nGo with %v goroutines and GOMAXPROCS=%v:\n", goroutines,
		os.Getenv("GOMAXPROCS"))
	fmt.Printf("Sum for n=1..%v of 1/n =\n", n)

	wg.Add(goroutines)
	partialSums = make([]float64, goroutines)
	rangeLen := n / goroutines
	for i := 0; i < goroutines; i++ {
		start := i*rangeLen + 1
		end := start + rangeLen - 1
		go calcPartialSum(i, start, end)
	}
	wg.Wait()
	sum := 0.0
	for i := 0; i < goroutines; i++ {
		sum += partialSums[i]
	}
	fmt.Printf("%v", sum)
}
