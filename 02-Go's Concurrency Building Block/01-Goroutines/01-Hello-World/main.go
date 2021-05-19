package main

import (
	"fmt"
	"sync"
)

var (
	wg sync.WaitGroup
)

func main() {
	sayHello := func() {
		defer wg.Done()
		fmt.Println("Hello world")
	}

	wg.Add(1)
	go sayHello()
	// this is the join point
	wg.Wait()
}
