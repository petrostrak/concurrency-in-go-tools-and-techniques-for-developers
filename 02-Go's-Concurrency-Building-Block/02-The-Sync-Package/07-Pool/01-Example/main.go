package main

import (
	"fmt"
	"sync"
)

func main() {
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating a new instance")
			return struct{}{}
		},
	}

	// Here we call Get on the pool. These calls will invoke the New function defined on
	// the pool since instances haven't yet been instantiated
	myPool.Get()
	instance := myPool.Get()

	// Here we put an instance previously retrieved back in the pool. This increases the available
	// number of instances to one.
	myPool.Put(instance)

	// When this call is executed, we will reuse the instance previously allocated and put it back
	// in the pool. The New function will not be invoked.
	myPool.Get()
}
