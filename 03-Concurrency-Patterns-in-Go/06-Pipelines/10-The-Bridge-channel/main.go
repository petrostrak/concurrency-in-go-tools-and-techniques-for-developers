package main

import "fmt"

// In some circumstances, you may find yourself wanting to consume values from a
// sequence of channels
// <-chan <-chan interface{}
//
// A sequence of channels suggests an ordered write, but from different sources.
func main() {
	orDone := func(done, c <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case valStream <- v:
					case <-done:
					}
				}
			}
		}()

		return valStream
	}

	// As a consumer, the code may be not care about the fact that its values come from a
	// sequence of channels. If we instead define a function that can destructure the channel
	// of channels into a simple channel -a technique called bridging the channels- this will
	// make it much easier for the consumer to focus on the problem at hand.
	bridge := func(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {

		// This is the channel that will return all the values from bridge.
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)

			// This loop is responsible for pulling channels off of chanStream and providing them to
			// a nested loop for use.
			for {
				var stream <-chan interface{}
				select {
				case maybeStream, ok := <-chanStream:
					if !ok {
						return
					}
					stream = maybeStream
				case <-done:
					return
				}

				// This loop is responsible for reading values off the channel it has been given and
				// repeating those values onto valStream. When the stream we're currently looping over
				// is closed, we break out of the loop performing the reads from this channel, and
				// continue with the next iteration of the loop, selecting channels to read from. This
				// provides us with an unbroken stream of values.
				for val := range orDone(done, stream) {
					select {
					case valStream <- val:
					case <-done:
					}
				}
			}
		}()

		return valStream
	}

	// Here's an example that creates a series of 10 channels, each with one element written to them,
	// and passes these channels into the bridge function.
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()

		return chanStream
	}

	for v := range bridge(nil, genVals()) {
		fmt.Printf("%v", v)
	}
}
