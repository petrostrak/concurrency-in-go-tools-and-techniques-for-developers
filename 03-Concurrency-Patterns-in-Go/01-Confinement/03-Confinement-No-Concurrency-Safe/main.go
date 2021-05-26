package main

import (
	"bytes"
	"fmt"
	"sync"
)

var (
	wg   sync.WaitGroup
	buff bytes.Buffer
)

func main() {
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}

		fmt.Println(buff.String())
	}

	wg.Add(2)
	data := []byte("golang")

	// Here we pass in a slice containing the first three bytes in the data structure
	go printData(&wg, data[:3]) // 3rd included

	// Here we pass in a slice containing the last three bytes in the data structure
	go printData(&wg, data[3:]) // 3rd excluded
	wg.Wait()
}
