package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	wg         sync.WaitGroup
	sharedLock sync.Mutex
)

const (
	runtime = 1 * time.Second
)

func main() {

	// greedyWorker greedily holds onto the shared lock for the entirety of its work loop
	greedyWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(3 * time.Nanosecond)
			sharedLock.Unlock()

			count++
		}

		fmt.Printf("Greedy worker was able to execute %v work loops\n", count)
	}

	// whereas the politeWorker attempts to only lock when it needs to. Both workers do
	// the same amount of simulated work, but in the same amount of time, the greedyWorker
	// got almost twice the amount of work done.
	//
	// The greedyWorker has unnecessrily expanded its hold on the shared lock beyond its critical
	// section and is preventing (via starvation) the polite worker's goroutine from performing work efficiently.
	politeWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			count++
		}

		fmt.Printf("Polite worker was able to execute %v work loops\n", count)
	}

	wg.Add(2)
	go greedyWorker()
	go politeWorker()
	wg.Wait()
}
