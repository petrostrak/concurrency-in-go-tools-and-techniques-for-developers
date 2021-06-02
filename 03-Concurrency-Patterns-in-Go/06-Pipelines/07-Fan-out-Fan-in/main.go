package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// This naive implementation of the fan out, fan in algorithm only works if the order in which results arrive is
// unimportant.
func main() {
	repeatFn := func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}

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

	toInt := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case intStream <- v.(int):
				}
			}
		}()
		return intStream
	}

	primeFinder := func(done <-chan interface{}, intStream <-chan int) <-chan interface{} {
		primeStream := make(chan interface{})
		go func() {
			defer close(primeStream)
			for integer := range intStream {
				integer -= 1
				prime := true
				for divisor := integer - 1; divisor > 1; divisor-- {
					if integer%divisor == 0 {
						prime = false
						break
					}
				}

				if prime {
					select {
					case <-done:
						return
					case primeStream <- integer:
					}
				}
			}
		}()
		return primeStream
	}

	// Fanning in means multiplexing of joining together multiple streams of data into a single stream.
	//
	// Here we take in our standard done channel to allow our goroutines to be torn down, and
	// then a variadic slice of interface{} channels to fan-in.
	fanIn := func(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {

		// On this line we create a sync.WaitGroup so that we can wait until all channels have been
		// drained.
		var wg sync.WaitGroup
		multiplexedStream := make(chan interface{})

		// Here we create a function, multiplex, which, when passed a channel, will read from the channel
		// and pass the value read onto the multiplexedStream channel.
		multiplex := func(c <-chan interface{}) {
			defer wg.Done()
			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		// Select from all the channels
		// This line increments the sync.WaitGroup by the number of channels we're multiplexing.
		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}

		// Wait for all the reads to complete
		// Here we create a goroutine to wait for all the channels we're multiplexing to be drained so
		// that we can close the multiplexedStream channel.
		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	rand := func() interface{} { return rand.Intn(50000000) }

	randIntStream := toInt(done, repeatFn(done, rand))

	// The process of fanning out a stage in a pipeline is exraordinarily easy. All we have to do
	// is to start multiple versions of that stage.
	// Here we're starting up as many copies of this stage as we have CPUs.
	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	finders := make([]<-chan interface{}, numFinders)
	fmt.Println("Primes:")
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}

	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v", time.Since(start))
}
