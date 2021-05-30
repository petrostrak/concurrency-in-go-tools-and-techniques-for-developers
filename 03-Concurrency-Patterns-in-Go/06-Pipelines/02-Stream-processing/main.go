package main

import "fmt"

// Stream processing means that the stage receives and emits one
// element at a time

func main() {
	multiply := func(value int, multiplier int) int {
		return value * multiplier
	}

	add := func(value int, additive int) int {
		return value + additive
	}

	integers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	// Each stage is receiving and emitting a discrete value, and the memory footprint
	// of our program is back down to only the size of the pipeline's input. But we had
	// to pull the pipeline down into the body of the for-loop and let the range do the
	// heavy lifting of feeding our pipeline. Not only does this limit the reuse of how
	// we feed the pipeline, but it also limits our ability to scale.
	//
	// Efectively, we're instantiating our pipeline for every iteration of the loop.
	// Thought it's cheap to make function calls, we're making three function calls for
	// each iteration of the loop.
	for _, v := range integers {
		fmt.Println(multiply(add(multiply(v, 2), 1), 2))
	}
}
