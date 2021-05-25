package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	// Here we create an in-memory buffer to help mitigate the nondeterministic nature
	// of the output. It doesn't give us any guarantees, but it's a little faster than
	// writing to stdout directly.
	var stdoutBuff bytes.Buffer

	// Here we ensure that the buffer is written out to stdout before the process exits.
	defer stdoutBuff.WriteTo(os.Stdout)

	// Here we create a buffered channel with the capacity of 4.
	intStream := make(chan int, 4)
	go func() {
		defer close(intStream)
		defer fmt.Fprintf(&stdoutBuff, "Producer Done.\n")

		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intStream <- i
		}
	}()

	for integer := range intStream {
		fmt.Fprintf(&stdoutBuff, "Received: %d\n", integer)
	}
}
