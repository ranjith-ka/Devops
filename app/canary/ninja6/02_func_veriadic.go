package main

import "fmt"

// func (receiver) identifier(parameters) (returns) { code }

func main() {
	x := foo12(2, 3, 4, 5, 6, 7, 8, 9)
	fmt.Println("The total is", x)

	// composite literals
	y := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	n := foo12(y...) // unfurling

	fmt.Println("The total is", n)
}

func foo12(s ...int) int {
	fmt.Println(s)
	sum := 0
	for _, v := range s {
		sum += v
		fmt.Println(sum)
	}
	return sum
}
