package main

import (
	"fmt"
	"sync"
)

var (
	once       sync.Once
	increments sync.WaitGroup
	count      int
)

func increment() {
	count++
}

func decrement() {
	count--
}

func countHundred() {
	increments.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}

	increments.Wait()
	fmt.Printf("Count is %d\n", count)
}

func countOnce() {
	once.Do(increment)
	once.Do(decrement)

	fmt.Printf("Count is %d\n", count)
}

func main() {
	// As the name implies, sync.Once is a type that utilizes some sync primitives internally
	// to ensure that only one call to Do ecer calls the function passed in - even on different goroutines
	countHundred()

	// It is surprising that the output displays 1 and not 0? This is because sync.Once only counts
	// the number of times Do is called, not how many times unique functions passed into Do are called
	countOnce()
}
