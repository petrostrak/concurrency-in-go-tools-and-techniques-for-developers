package main

import "fmt"

func main() {
	intStream := make(chan int)

	go func() {
		// Here we ensure that the channel is closed before we exit the goroutine. This
		// is a very common pattern
		defer close(intStream)
		for i := 1; i < 5; i++ {
			intStream <- i
		}
	}()

	// Here we range over the intStream
	for integer := range intStream {
		fmt.Printf("%v ", integer)
	}
}
