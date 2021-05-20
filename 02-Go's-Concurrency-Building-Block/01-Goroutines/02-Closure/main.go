package main

import (
	"fmt"
	"sync"
)

var (
	wg sync.WaitGroup
)

func main() {
	salutation := "hello"

	wg.Add(1)
	go func() {
		defer wg.Done()
		// the goroutine will modify the value of the variable salutation
		salutation = "welcome"
	}()

	wg.Wait()
	fmt.Println(salutation)
}
