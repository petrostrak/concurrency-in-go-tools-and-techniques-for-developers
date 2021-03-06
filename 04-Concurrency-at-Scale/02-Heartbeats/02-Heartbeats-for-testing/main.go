package main

import (
	"fmt"
	"math/rand"
)

func main() {
	doWork := func(done <-chan interface{}) (<-chan interface{}, <-chan int) {

		// Here we create the heartbeat channel with a buffer of one. This ensures that
		// there's always at least one pulse sent tout even if no one is listening in
		// time for the send to occur.
		heartbeatStream := make(chan interface{}, 1)
		workStream := make(chan int)
		go func() {
			defer close(heartbeatStream)
			defer close(workStream)

			for i := 0; i < 10; i++ {

				// Here we set up a separate select block for the heartbeat. We don't want to
				// include this in the same select block as the send on results because if the
				// receiver isn't ready for the result, they'll receive a pulse instead, and
				// the current value of the result will be lost. We also don't include a case
				// statement for the done channel since we have a default case that will just
				// fall through.
				select {
				case heartbeatStream <- struct{}{}:

				// Once again we guard against the fact that no one may be listening to our
				// heartbeats. Because our heartbeat channel was created with a buffer of one,
				// if someone is listening, but not in time for the first pulse, they'll still
				// be notified of a pulse.
				default:
				}

				select {
				case <-done:
					return
				case workStream <- rand.Intn(10):
				}
			}
		}()

		return heartbeatStream, workStream
	}

	done := make(chan interface{})
	defer close(done)

	heartbeat, results := doWork(done)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)
			} else {
				return
			}
		}
	}
}
