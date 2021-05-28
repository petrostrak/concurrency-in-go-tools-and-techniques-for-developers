package main

import (
	"fmt"
	"time"
)

var (
	or func(channels ...<-chan interface{}) <-chan interface{}
)

func main() {

	// Here we have our function, or, which takes in a variadic slice of channels and
	// returns a single channel.
	or := func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {

		// Since this is a recursive function, we must set up termination criteria. The
		// first is that if the variadic slice is empty, we simply return a nil channel.
		//This is consistant with the idea of passing in no channels; we wouldn't expect
		// a composite channel to do anything.
		case 0:
			return nil

		// Our second determination criteria states that if our variadic slice only contains
		// one element, we just return that element.
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})

		// Here is the main body of the function, and where the recursion happens. We create a
		// goroutine so that we can wait for messages on our channels without blocking.
		go func() {
			defer close(orDone)

			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}

			// Because of how we're recursing, every recursive call to or will at least have
			// two channels. As an optimization to keep the number of goroutines constrained,
			// we place a special case here for calls to or with only two channels.
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:

				// Here we recursively create an or-channel from all the channels in our slice
				// after the third index, and then select from this. This recurrence relation
				// will destructure the rest of the slice into or-channels to form a tree from
				// which the first signal will return. We also pass on the orDone channel so that
				// when goroutines up the tree exit, goroutines down the tree also exit.
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()

		return orDone
	}

	// This function simply creates a channel that will close when the time specified in
	// the after elapses.
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()

		return c
	}

	// Here we keep track of roughly when the channel from the or function begins to block
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	// And here we print the time it took for the read to occur.
	fmt.Printf("Done after %v", start)
}
