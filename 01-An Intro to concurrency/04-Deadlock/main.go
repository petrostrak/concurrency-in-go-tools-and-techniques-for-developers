package main

import (
	"fmt"
	"sync"
	"time"
)

// A deadlocked program is one in which all concurrent processes are waiting on one another.
// In this state, the program will never recover without outside intervention.

type value struct {
	mu    sync.Mutex
	value int
}

var (
	wg sync.WaitGroup
)

// fatal error: all goroutines are asleep - deadlock!
func main() {
	printSum := func(v1, v2 *value) {
		defer wg.Done()
		// we attempt to enter the critical section for the incoming value
		v1.mu.Lock()
		// we use defer to exit the critial section before printSum returns
		defer v1.mu.Unlock()

		// we sleep to simulate work(and trigger a deadlock)
		time.Sleep(2 * time.Second)
		v2.mu.Lock()
		defer v2.mu.Unlock()

		fmt.Printf("sum = %v\n", v1.value+v2.value)
	}

	var a, b value
	wg.Add(2)
	// first call to printSum locks a and then attempts to lock b
	go printSum(&a, &b)
	// in the meantime second call to printSum has locked b and attempts to lock a.
	// Both goroutines wait infinitely on each other
	go printSum(&b, &a)
	wg.Wait()
}
