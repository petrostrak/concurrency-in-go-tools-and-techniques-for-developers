package main

import (
	"fmt"
	"time"
)

func main() {
	doWork := func(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {

		// Here we set up a channel to send heartbeats on. We return this out of doWork.
		heartBeat := make(chan interface{})
		results := make(chan time.Time)
		go func() {
			defer close(heartBeat)
			defer close(results)

			// Here we set the heartbeat to pulse at the pulseInterval we were given. Every pulseInterval
			// there will be something to read on this channel.
			pulse := time.Tick(pulseInterval)

			// This is just another thicker used to simulate work coming in. We choose a duration greater
			// than the pulseInterval so that we can see some heartbeats coming out of the goroutine.
			workGen := time.Tick(2 * pulseInterval)

			sendPulse := func() {
				select {
				case heartBeat <- struct{}{}:

				// Note that we include a default clause. We must always guard against the fack that no one
				// may be listening to our heartbeat. The results emitted from the goroutine are critical,
				// but the pulses are not.
				default:
				}
			}

			sendResult := func(r time.Time) {
				for {
					select {
					case <-done:
						return

					// Just like with done channels, anytime you perform a send or receive, you also need to
					// include a case for the heartbeat's pulse.
					case <-pulse:
						sendPulse()
					case results <- r:
						return
					}
				}
			}

			for {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()

		return heartBeat, results
	}

	done := make(chan interface{})

	// We set up the standard done channel and close it after 10 seconds. This gives our
	// goroutine time to do some work.
	time.AfterFunc(10*time.Second, func() { close(done) })

	// Here we set our timeout period. We'll use this to couple our heartbeat interval
	// to our timeout.
	const timeout = 2 * time.Second

	// We pass in timeout/2 here. This gives our heartbeat an extra tick to respond so
	// that our timeout isn't too sensitive.
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {

		// Here we select on the heartbeat. When there are no results, we are at least
		// guaranteed a message from the heartbeat channel every timeout/2. If we don't
		// receive it, we know there's something wrong with the goroutine itself.
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("pulse")

		// Here we select from the results channel.
		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Printf("results %v\n", r.Second())

		// Here we time out if we haven't received either a heartbeat or a new result.
		case <-time.After(timeout):
			return
		}
	}
}
