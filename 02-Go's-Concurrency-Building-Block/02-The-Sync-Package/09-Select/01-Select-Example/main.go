package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()

	c := make(chan interface{})

	go func() {
		time.Sleep(5 * time.Second)

		// Here we close the channel after waiting 5 seconds.
		close(c)
	}()

	fmt.Println("Blocking on read...")

	select {

	// Here we attempt a read on the channel. Note that as this code is written, we
	// don't requre a select statement - we could simply write <-c.
	case <-c:
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}
}
