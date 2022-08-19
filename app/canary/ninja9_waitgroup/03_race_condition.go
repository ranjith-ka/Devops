package main

import (
	"fmt"
	"sync"
)

var (
	counter int32          // counter is a variable incremented by all goroutines.
	wg      sync.WaitGroup // wg is used to wait for the program to finish.
	mutex   sync.Mutex     // mutex is used to define a critical section of code.
)

func main() {
	wg.Add(3) // Add a count of two, one for each goroutine.

	go increment("Realme")
	go increment("Oppo")
	go increment("Apple")

	wg.Wait() // Wait for the goroutines to finish.
	fmt.Println("Counter:", counter)

}

func increment(lang string) {
	defer wg.Done() // Schedule the call to Done to tell main we are done.
	for i := 0; i < 3; i++ {
		mutex.Lock()
		{
			fmt.Println(lang)
			counter++
		}
		mutex.Unlock()
	}
}
