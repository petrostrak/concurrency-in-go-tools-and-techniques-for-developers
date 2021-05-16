package main

import "fmt"

func main() {
	var i int

	go func() {
		i++
	}()

	if i == 0 {
		fmt.Printf("The value of i is %d\n", i)
	}
}
