package main

import (
	"fmt"
	"time"
)

var (
	c1, c2 <-chan int
)

func main() {
	start := time.Now()

	select {
	case <-c1:
	case <-c2:
	default:
		fmt.Printf("In default after %v\n", time.Since(start))
	}
}
