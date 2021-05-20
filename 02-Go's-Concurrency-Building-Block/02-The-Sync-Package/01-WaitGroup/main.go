package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

func main() {

	// We can Add with an argument of 1 to indicate that one goroutine is
	// begining.
	wg.Add(1)
	go func() {
		// Here we call Done using the defer keywork to ensure that before
		// we exit the goroutine's closure, we indicate to the WaitGroup
		// that we have exited.
		defer wg.Done()
		fmt.Println("1st goroutine sleeping...")
		time.Sleep(1)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("2nd goroutine sleeping...")
		time.Sleep(2)
	}()

	// Here we call Wait which will block the main goroutine, until all goroutines
	// have indicated that they have exited.
	wg.Wait()
	fmt.Println("All goroutines complete.")
}
