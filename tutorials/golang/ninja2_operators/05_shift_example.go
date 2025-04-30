package main

import (
	"fmt"
)

func main() {
	var t, i uint
	t, i = 1, 1

	for i = 1; i < 10; i++ {
		fmt.Printf("%d << %d = %d \n", t, i, t<<i) // left shifting increment
	}

	fmt.Println()

	t = 512
	for i = 1; i < 10; i++ {
		fmt.Printf("%d >> %d = %d \n", t, i, t>>i) // Decrement for right
	}
}
