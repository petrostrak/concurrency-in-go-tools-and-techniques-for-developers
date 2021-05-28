package main

import (
	"fmt"
	"net/http"
)

// Result is a type that encompasses both the *http.Response and error possible
// from an iteration of the loop within our goroutine.
type Result struct {
	Error    error
	Response *http.Response
}

func main() {

	// This line returns a channel that can be read from to retrieve results of an iteration
	// of our loop.
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, url := range urls {

				// Here we create a Result instance with Error and Response fields set.
				var result Result
				resp, err := http.Get(url)
				result = Result{Error: err, Response: resp}
				select {
				case <-done:
					return

				// This is where we write the Result to the channel.
				case results <- result:
				}
			}
		}()

		return results
	}

	done := make(chan interface{})
	defer close(done)

	errCount := 0

	urls := []string{"https://www.google.com", "a", "b", "c", "d"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}

}
