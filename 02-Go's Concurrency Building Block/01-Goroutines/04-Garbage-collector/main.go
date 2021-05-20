package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var c <-chan interface{}
	var wg sync.WaitGroup

	// we require a goroutine that will never exit so that we can keep a number of them
	// in memory for measurement. The goroutine won't exit until the process is finished.
	noop := func() {
		wg.Done()
		<-c
	}

	// we define the number of goroutines to create.
	// we use the law of large numbers to asymptotically approach
	// the size of a goroutine.
	const numGoroutines = 1e4
	wg.Add(numGoroutines)

	// we measure the amount of memory consumed before creating our goroutines.
	before := memConsumed()

	for i := numGoroutines; i > 0; i-- {
		go noop()
	}

	wg.Wait()

	// we measure the amount of memory consumed after creating our goroutines
	after := memConsumed()
	fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
}
