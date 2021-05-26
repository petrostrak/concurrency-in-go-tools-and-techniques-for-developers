package main

// // Either loop infinitely or range over something
// for {
// 	select {
// 		// Do some work with channels
// 	}
// }

var (
	done, stringSteam chan interface{}
)

func main() {

	for _, s := range []string{"a", "b", "c"} {
		select {
		case <-done:
			return
		case stringSteam <- s:
		}
	}
}
