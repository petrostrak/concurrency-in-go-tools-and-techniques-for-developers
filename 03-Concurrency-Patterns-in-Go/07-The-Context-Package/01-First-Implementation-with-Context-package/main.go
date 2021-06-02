package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

func main() {

	// Here main creates a new Context with context.Background() and wraps it with
	// context.WithCancel to allow for cancellations.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printGreeting(ctx); err != nil {
			fmt.Printf("cannot print greeting: %v\n", err)

			// (*)
			// On this line, main will cancel the Context if there is an error returned from printGreeting.
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printFarewell(ctx); err != nil {
			fmt.Printf("cannot print farewell: %v", err)

			// (*)
			// On this line, main will cancel the Context if there is an error returned from printFarewell.
			cancel()
		}
	}()

	wg.Wait()

}

func genGreeting(ctx context.Context) (string, error) {

	// Here genGreeting wraps its Context with context.WithTimeout. This will automatically
	// cancel the returned Context after 1 second, thereby canceling any children it passes
	// the Context into, namely locale.
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	switch locale, err := locale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}

	return "", fmt.Errorf("unsupported locale")
}

func genFarewell(ctx context.Context) (string, error) {

	// Here genFarewell wraps its Context with context.WithTimeout. This will automatically
	// cancel the returned Context after 1 second, thereby canceling any children it passes
	// the Context into, namely locale.
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	switch locale, err := locale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}

	return "", fmt.Errorf("unsupported locale")
}

func locale(ctx context.Context) (string, error) {

	// Here we check to see whether our Context has provided a deadline. If it did,
	// and our system's clock has advanced past the deadline, we simply return with
	// a special error defined in the context package, DeadlineExceeded.
	if deadline, ok := ctx.Deadline(); ok {
		if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0 {
			return "", context.DeadlineExceeded
		}
	}
	select {
	case <-ctx.Done():

		// This line returns the reason why the Context was canceled. This error will bubble
		// all the way up to main, which cause the cancellation at main(*).
		return "", ctx.Err()
	case <-time.After(1 * time.Minute):
	}

	return "EN/US", nil
}

func printGreeting(ctx context.Context) error {
	greeting, err := genGreeting(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("%s world!\n", greeting)
	return nil
}

func printFarewell(ctx context.Context) error {
	farewell, err := genFarewell(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("%s world!\n", farewell)
	return nil
}
