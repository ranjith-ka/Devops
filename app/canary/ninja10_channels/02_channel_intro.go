package main

import (
	"fmt"
)

func nija10channels() {
	numCh := make(chan int)
	done := make(chan struct{})
	go counter(numCh)
	go printer(numCh, done)


func counter(out chan int) {
	for i := 0; i < 10; i++ {
		out <- i
	}
	close(out)
}

func printer(numCh chan int) {
	for n := range numCh {
		fmt.Println(n)
		}	
	}
}
