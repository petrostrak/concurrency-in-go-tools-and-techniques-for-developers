package main

import (
	"fmt"
	"time"
)

func main() {
	// Here we pass the done channel to the doWork function. As a convention, this channel
	// is the first parameter.
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		termindated := make(chan interface{})

		go func() {
			defer fmt.Println("doWork exited.")
			defer close(termindated)

			for {
				select {
				case s := <-strings:
					// Do something interesting
					fmt.Println(s)

				// On this line we see the ubiquitous for-select pattern in use. One of our case statements
				// is checking whether our done channel has been signaled. If it has, we return from the
				// goroutine.
				case <-done:
					return
				}
			}
		}()

		return termindated
	}

	done := make(chan interface{})
	terminated := doWork(done, nil)

	// Here we create another goroutine that will cancel the goroutine spawned in doWork if more than
	// one second passes.
	go func() {
		// Cancel the operation after 1 second.
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	// This is where we join the goroutine spawned from doWork with the main goroutine.
	<-terminated
	fmt.Println(done)
}
