// passing the function as an Argument to other function

// callback

package main

import (
	"fmt"
)

func main() {
	ii := []int{1, 2, 3, 4, 5, 6, 7, 8}
	s := sum(ii...)
	fmt.Println(s)
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

// function for call back

// func even(f func(xi ...int) int, vi ...int) int {
//     yi := []int
// 	return yi
// }
