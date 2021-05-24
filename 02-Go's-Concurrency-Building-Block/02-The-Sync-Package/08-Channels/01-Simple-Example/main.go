package main

import "fmt"

func main() {
	stringStream := make(chan string)

	go func() {
		// Here we pass a string literal onto the channel stringStream
		stringStream <- "Hello channels!"
	}()

	// Here we read the string literal off of the channel and print it out to stdout
	fmt.Println(<-stringStream)
}

// Why the anonymous goroutine completes before the main goroutine?
//
// This example works because channels in Go are said to be blocking. This means that
// any goroutine that attempts to write to a channel that is full, will wait until the
// channel has been emptied, and any goroutine that attempts to read from a channel that
// is empty, will wait until at least one item is placed on it.
// Thus, the main goroutine and the anonymous goroutine block deterministically.
