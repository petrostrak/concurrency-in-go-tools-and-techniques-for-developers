package main

import (
	"fmt"
	"time"
)

var (
	c <-chan int
)

func main() {
	select {

	// This case statement will never become unblocked because we're reading from
	// a nil channel
	case <-c:
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out")
	}
}
