package main

import (
	"sync"
	"testing"
)

// go test -bench=. -cpu-1 main_test.go
func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})

	var token struct{}

	sender := func() {
		defer wg.Done()
		// here we wait until we're told to begin. We don't want the cost of setting up and
		// starting each gorutine to factor into the measurement of context switching
		<-begin

		for i := 0; i < b.N; i++ {
			// here we send messages to the receiver goroutine. A struct{}{} is called an empty
			// struct and takes up no memory, thus we only measuring the time it takes to signal
			// a message
			c <- token
		}
	}

	receiver := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			// here we receive a message but do nothing with it
			<-c
		}
	}

	wg.Add(2)
	go sender()
	go receiver()
	// here we start the performance timer
	b.StartTimer()
	// here we tell the two goroutines to begin
	close(begin)
	wg.Wait()
}
