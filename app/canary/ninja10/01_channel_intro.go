package main

import (
	"fmt"
)

func main() {
	// c := make(chan int)      // Create a deafult channel to write

	// c := make(chan<- int)  // Send only channel

	// c := make(<-chan int)  // Recive only channel

	c := make(chan int, 1) // Buffer channels
	c <- 42                // Channel Blocks, read only or send only channel , open channel to write but close if you want to read, else use buffer.
	// go func() {
	// 	c <- 42
	// }()

	fmt.Println(<-c) // Recive from channel

	fmt.Printf("------\n")
	fmt.Printf("c\t%T\n", c)
}
