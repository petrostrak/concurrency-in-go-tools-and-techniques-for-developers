package main

import "fmt"

// A pipeline is nothing more that a series of things that take data in,
// perform an operation on it, and pass the data back out. Each of these
// operations are called stages of the pipeline.
//
// By using pipelines, we separate the concerns of each stage. We can modify
// stages independent of one another, we can mix and match how stages are
// combined independent of midifying the stages, we can process each stage
// concurrent to upstream or downstream stages, and we can fan-out or rate-limit
// portions of our pipeline.
//
// The properties of a pipeline stage are:
//
// A stage consumes and returns the same type.
// A stage must be reified by the language so that it may be passed around. Functions
// in Go are reified and fit this purpose nicely.

func main() {

	multiply := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, v := range values {
			multipliedValues[i] = v * multiplier
		}

		return multipliedValues
	}

	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}

		return addedValues
	}

	integers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	// Here, our add and multiply stages satisfy all the properties of a pipeline stage:
	// They both consume a slice of int and return a slice of int, and because Go has
	// reified functions, we can pass add and multiply around.
	for _, v := range add(multiply(integers, 2), 1) {
		fmt.Println(v)
	}

	// If we wanted to now add an additional stage to our pipeline to multiply by two, we'd
	// simply wrap our previous pipeline in a new multiply stage like so:
	for _, v := range multiply(add(multiply(integers, 2), 1), 2) {
		fmt.Println(v)
	}
}
