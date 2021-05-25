package main

import "fmt"

func main() {

	ownerChan := func() <-chan int {
		// Here we instantiate a buffered channel. Since we know we'll produce 6 results,
		// we can create a buffered channel of five so that the goroutine can complete as
		// quickly as possible.
		resultStream := make(chan int, 5)

		// Here we start an anonymous goroutine that performs writes on resultStream. The
		// goroutine is now encapsulated within the surrounding function.
		go func() {

			// Here we ensure that resultStream is closed once we're finished with it. As the
			// channel owner, this is our responsibility.
			defer close(resultStream)

			for i := 0; i <= 5; i++ {
				resultStream <- i
			}
		}()

		// Here we return the resultStream channel. Since the return value is declared as a
		// read-only channel, resultStream will implicitly be converted to read-only for
		// consumer.
		return resultStream
	}

	resultStream := ownerChan()

	// Here we range over resultStream. As a consumer, we are only concerned with blocking and
	// closed channels.
	for result := range resultStream {
		fmt.Printf("Received: %d\n", result)
	}

	fmt.Println("Done receiving")
}
