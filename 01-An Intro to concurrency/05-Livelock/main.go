package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	cadance := sync.NewCond(&sync.Mutex{})

	go func() {
		for range time.Tick(1 * time.Millisecond) {
			cadance.Broadcast()
		}
	}()

	takeStep := func() {
		cadance.L.Lock()
		cadance.Wait()
		cadance.L.Unlock()
	}

	// tryDir allows a person to attempt to move in a direction and returns
	// whether or not they were successful. Each direction is represented as a
	// count of the number of people trying to move in that direction.
	tryDir := func(dirName string, dir *int32, out *bytes.Buffer) bool {
		fmt.Fprintf(out, " %v", dirName)

		// we declare out intention to move in a direction by incrementing that
		// direction by one. Operation is atomic.
		atomic.AddInt32(dir, 1)

		// to demonstrate a livelock, each person must move at the same rate of speed,
		// or cadence. takeStep() simulates a constant cadence between all parties.
		takeStep()
		if atomic.LoadInt32(dir) == 1 {
			fmt.Fprintf(out, ". Success!")
			return true
		}

		takeStep()

		// the person realizes they cannot go in this direction and gives up. We indicate
		// this by decrementiong that direction by one.
		atomic.AddInt32(dir, -1)
		return false
	}

	var left, right int32

	tryLeft := func(out *bytes.Buffer) bool {
		return tryDir("left", &left, out)
	}

	tryRight := func(out *bytes.Buffer) bool {
		return tryDir("right", &right, out)
	}

	walk := func(walking *sync.WaitGroup, name string) {
		var out bytes.Buffer
		defer func() {
			fmt.Println(out.String())
		}()
		defer walking.Done()

		fmt.Fprintf(&out, "%v is trying to scoot:", name)

		// artificial limit on the number of attempts
		for i := 0; i < 5; i++ {

			// the person will attempt to step left, and if that fails,
			// they will attempt to step right
			if tryLeft(&out) || tryRight(&out) {
				return
			}
		}

		fmt.Fprintf(&out, "\n%v tosses her hands up in exasperation!", name)
	}

	// the variable provides a way for the programm to wait until both people are either
	// able to pass one another, or give up
	var peopleInHallway sync.WaitGroup
	peopleInHallway.Add(2)
	go walk(&peopleInHallway, "Alice")
	go walk(&peopleInHallway, "Barbara")
	peopleInHallway.Wait()
}
