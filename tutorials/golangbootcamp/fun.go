package main

import (
	"fmt"
)

func branch(size, ratio, angle float64, iterations int) {
	// draw the stem
	fmt.Println("color black")
	fmt.Printf("width %f\n", size)
	fmt.Printf("forward %f\n", size)
	fmt.Println("color off")

	if iterations > 0 {
		// rotate left and draw the child stem
		fmt.Printf("left %f\n", angle)
		branch(size*ratio, ratio, angle, iterations-1)

		// rotate right and draw the child stem
		fmt.Printf("right %f\n", angle*2)
		branch(size*ratio, ratio, angle, iterations-1)

		// rotate back to the main stem position
		fmt.Printf("left %f\n", angle)
	}

	// rotate back and walk to the start of the stem
	fmt.Println("left 180")
	fmt.Printf("forward %f\n", size)

	// rotate to the initial position
	fmt.Println("left 180")
}

func main() {
	fmt.Println("draw mode")
	branch(5, 0.7, 30, 5)
}
