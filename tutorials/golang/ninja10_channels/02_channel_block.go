package main

import "fmt"

func main() {
	c := gen()
	receive(c)

	fmt.Println("about to exit")
}

func receive(c <-chan int) {
	for v := range c {
		fmt.Println(v) // Blocks until we read the channel
	}
}

func gen() <-chan int { // func return channel
	c := make(chan int)

	go func() { // New go routine to start concurrency. Channel to sync between go routine.
		for i := 0; i < 10; i++ {
			c <- i // Also Blocks until receive concludes
		}
		close(c)
	}()
	return c
}
