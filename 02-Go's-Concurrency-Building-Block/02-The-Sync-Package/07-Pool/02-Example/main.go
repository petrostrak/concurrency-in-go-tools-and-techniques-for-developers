package main

import (
	"fmt"
	"sync"
)

var (
	numCalcsCreated int
	wg              sync.WaitGroup
)

const (
	numWorkers = 1024 * 1024
)

func main() {
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024)

			// Notice that we are storing the address of the slice of bytes
			return &mem
		},
	}

	// Seed the pool with 4KB
	calcPool.Put(calcPool.New)
	calcPool.Put(calcPool.New)
	calcPool.Put(calcPool.New)
	calcPool.Put(calcPool.New)

	wg.Add(numWorkers)

	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()

			mem := calcPool.Get()
			defer calcPool.Put(mem)

			// Assume something interesting, but quick is being done
			// with this memory
		}()
	}

	wg.Wait()
	fmt.Printf("%d calculators were created.\n", numCalcsCreated)
}
