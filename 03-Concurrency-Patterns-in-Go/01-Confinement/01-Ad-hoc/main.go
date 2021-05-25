package main

import "fmt"

func main() {

	// Data slice of integers is available from both the loopData function and the
	// loop over the handleData. This way, mistakes might be made, and the confinement
	// might break down and cause issues.
	data := make([]int, 4)

	loopData := func(handleData chan<- int) {
		defer close(handleData)

		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}
