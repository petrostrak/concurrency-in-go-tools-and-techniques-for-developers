package main

import (
	"fmt"
	"sync"
)

var (
	wg sync.WaitGroup
)

func main() {
	begin := make(chan interface{})

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// Here the goroutine waits until it is told it can continue
			<-begin
			fmt.Printf("%v has begun\n", i)
		}(i)
	}

	fmt.Println("Unblocking goroutines..")

	// Here we close the channel, thus unblocking all goroutines simultaneously
	close(begin)
	wg.Wait()
}
