package main

import "fmt"

func main() {

	// The generator function converts a discrete set of values into a stream of
	// data on a channel.
	//
	// The generator function takes in a variadic slice of integers..
	generator := func(done <-chan interface{}, integers ...int) <-chan int {

		// ..constructs a buffered channel of integers with a length equal to the
		// incoming integer slice..
		intStream := make(chan int)

		// ..starts a goroutine and returns a constructed channel.
		go func() {
			defer close(intStream)

			// Then, on the goroutine that was created, generator ranges over the
			// variadic slice that was passed in..
			for _, i := range integers {
				select {
				case <-done:
					return

				// ..and sends the slices values on the channel it created.
				case intStream <- i:
				}
			}
		}()

		return intStream
	}

	multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- multiplier * i:
				}
			}
		}()

		return multipliedStream
	}

	add := func(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case addedStream <- additive + i:
				}
			}
		}()

		return addedStream
	}

	// The first thing our program does is create a done channel and call close on it in a
	// defer statement. As discussed previously, this ensures our program exits cleanly and
	// never leaks goroutines.
	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4, 5, 6, 7, 8, 9)

	// For a stream of numbers, we'll multiply them by 2, add 1 and then multiply the result by 2.
	// This pipeline is similar to the 02-Stream-processing pipeline, but different in very important
	// ways.
	//
	// First, we are using channels. We can use a range statement to extract the values, and at each
	// stage we can safely execute concurrently, because our inputs and outputs are sage in concurrent
	// context.
	// Second, each stage of the pipeline is executing concurrently. This means that any stage only need
	// to wait for its inputs.
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	// Finally, we range over this pipeline and values are pulled through the system.
	for v := range pipeline {
		fmt.Println(v)
	}
}
