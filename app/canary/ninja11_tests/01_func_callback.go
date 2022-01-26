package main

import (
	"fmt"
)

func main() {

	// function expression to print total, assing the func to the variable
	TotalSum := func(xi ...int) int {
		total := 0
		for _, v := range xi {
			total += v
		}
		return total
	}
	ii := []int{1, 2, 3, 4, 5, 6, 7, 8}
	h := sum(ii...)
	fmt.Printf("Sum is %d\n", TotalSum(ii...))
	fmt.Printf("Sum is %d\n", h)
	operate()
}

func operate() {
	x := 0
	fmt.Println(x)
	x++
	fmt.Println(x)
	x += 220
	fmt.Println(x)
	x--
	fmt.Println(x)
	x -= 30
	fmt.Println(x)
}

// function to print total
func sum(xi ...int) int {
	total := 0
	for _, v := range xi {
		total += v
	}
	return total
}
