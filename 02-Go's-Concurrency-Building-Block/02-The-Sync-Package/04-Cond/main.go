package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// First we create our condition using a standard sync.Mutex as the Locker
	c := sync.NewCond(&sync.Mutex{})
	// Next, we create a slice with a length of zero. Sence we know we'll eventuallly
	// add 10 times, we instantiate it with a capacity of 10
	queue := make([]interface{}, 0, 10)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)

		// We once again enter the critical section for the condition so we can modify data perinent
		// to the condition
		c.L.Lock()

		// Here we simulate dequeuing an item by reassigning the head of the slice to the second item
		queue = queue[1:]
		fmt.Println("Removed from queue")

		// Here we exit the condition's critical section since we've successfully dequeued an item
		c.L.Unlock()

		// Here we let a goroutine waiting on the condition know that something has occured
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		// We enter the critical section for the condition by calling Lock on the condition's Locker
		c.L.Lock()
		// Here we check the length of the queue on a loop. This is important because a signal on the
		// condition doesn't necessarily mean that you've been waiting for has occured - only that
		// something has occured
		for len(queue) == 2 {
			// We call Wait, which will suspend the main goroutine until a signal on the conditin has been
			// sent
			c.Wait()
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})

		// Here we create a new goroutine that will dequeue an element after one second
		go removeFromQueue(1 * time.Second)

		// Here we exit the condition's critical section since we've successfully enqueued an
		// item
		c.L.Unlock()
	}
}
