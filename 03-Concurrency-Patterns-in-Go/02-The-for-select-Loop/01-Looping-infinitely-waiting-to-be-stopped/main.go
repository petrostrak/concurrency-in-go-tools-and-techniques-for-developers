package main

var (
	done, stringSteam chan interface{}
)

func main() {
	for {
		select {
		case <-done:
			return
		default:
		}

		// Do non-preemtable work
	}
}
