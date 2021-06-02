package main

import "fmt"

// Sometimes we may want to split values coming in from a channel so that we can
// send them off into two separate areas of our codebase.
//
// Taking its name from the tee command in Unix-like systems, the tee-channel does
// just this. You can pass it a channel to read from, and it will return two separate
// channels that will get the same value.
func main() {
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

	tee := func(done <-chan interface{}, in <-chan interface{}) (_, _ <-chan interface{}) {
		out1 := make(chan interface{})
		out2 := make(chan interface{})
		go func() {
			defer close(out1)
			defer close(out2)

			// Notice that writes to out1 and out2 are tightly coupled. The iteration over in
			// cannot continue until both out1 and out2 have been written to.
			for val := range orDone(done, in) {
				// We will want to use local versions of out1 and out2, so we shadow these
				// variables.
				var out1, out2 = out1, out2

				// We're going to use one select state ment so that writes to out1 and out2
				// don't block each other. To ensure both are written to, we'll perform two
				// iterations of the select statement: one for each outbound channel.
				for i := 0; i < 2; i++ {
					select {
					case <-done:

					// Once we've written to a channel, we set its shadowed copy to nil so that
					// further writes will block and the other channel may continue.
					case out1 <- val:
						out1 = nil
					case out2 <- val:
						out2 = nil
					}
				}
			}
		}()

		return out1, out2
	}

	done := make(chan interface{})
	defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 1, 2), 4))

	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
}
