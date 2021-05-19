package main

import (
	"fmt"
	"sync"
)

var (
	wg sync.WaitGroup
)

func main() {
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		// we declare a parameter, just like any other function. We shadow the
		// original salutation variable to make what's happening more apparent
		go func(salutation string) {
			defer wg.Done()
			fmt.Println(salutation)
		}(salutation) // here we pass in the current iteration's variable to the closure.
		// a copy of the string struct is make, thereby ensuring that when
		// the goroutine is run, we refer to the proper string
	}
	wg.Wait()
}
