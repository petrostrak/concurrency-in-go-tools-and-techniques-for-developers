package main

import (
	"fmt"
	"sync"
)

type Button struct {
	Clicked *sync.Cond
}

func main() {
	// We define a type Button that contains a condition, Clicked
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	// Here we define a coveniece function that will allow us to register functions to
	// handle signals from a condition. Each handler is run on its won goroutine, and
	// subscribe will not exit until that goroutine is confirmed to be running
	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)

		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()

		goroutineRunning.Wait()
	}

	// Here we set a handler for when the mouse button is raised. It in turn calls Broadcast on
	// the Clicked Cond to let all handlers to know that the mouse button has been clicked.
	var clickedRegistered sync.WaitGroup
	clickedRegistered.Add(3)

	// Here we create a WaitGroup. This is done only to ensure our program doesn't exit before
	// our writes to stdout occur
	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing window.")
		clickedRegistered.Done()
	})

	// Here we register a handler that simulates maximizing the button's window when the button is clicked
	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box")
		clickedRegistered.Done()
	})

	// Here we register a handler that simulates displaying a dialog box when the mouse is clicked
	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked")
		clickedRegistered.Done()
	})

	// Next, we simulate a user raising the mouse button from having clicked the application's button
	button.Clicked.Broadcast()

	clickedRegistered.Wait()
}
