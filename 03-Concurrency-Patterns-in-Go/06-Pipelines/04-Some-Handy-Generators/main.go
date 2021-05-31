package main

import "fmt"

// As a reminder, a generator for a pipeline is any function that converts a set of
// discrete values into a stream of values on a channel.

func main() {

	// This function will repeat the values you pass to it infinitely until you tell it
	// to stop.
	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()

		return valueStream
	}

	// This pipeline stage will only take the first num items off of its incoming valueStream and then
	// exit.
	take := func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()

		return takeStream
	}

	done := make(chan interface{})
	defer close(done)

	// We created a repeat generator to generate an infinite number of ones, but
	// then only take the first ten. Because the repeat generator's send blocks on
	// the take stage's receive, the repeat generator is very efficient. Although we
	// have the capability of generating an infinite stream of ones, we only generate
	// N+1 instances where N is the number we pass into the take stage.
	for num := range take(done, repeat(done, 1), 10) {
		fmt.Printf("%v ", num)
	}
}
