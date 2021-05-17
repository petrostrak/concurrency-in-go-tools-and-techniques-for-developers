package main

import (
	"fmt"
	"sync"
)

var (
	accessMemory sync.Mutex
)

func main() {
	var data int

	// We have 3 critical sections:
	// a) goroutine increments the data variable
	go func() {
		accessMemory.Lock()
		data++
		accessMemory.Unlock()
	}()

	accessMemory.Lock()
	// b) if statement checks the value of data
	if data == 0 {
		fmt.Println("The value is 0.")
	} else {
		// c) fmt.Printf statement retrieves the value of data
		fmt.Printf("The value is: %v. \n", data)
	}
	accessMemory.Unlock()
}
