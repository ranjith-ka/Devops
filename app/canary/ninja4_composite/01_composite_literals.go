package main

import (
	"fmt"
)

func main() {
	x := [5]int{1, 5, 3, 4, 2} // Slice

	// x := [5]int{0, 1, 4, 3, 2}  //  Array
	fmt.Println(x)
	for i, v := range x {
		fmt.Println(i, v)

	}
	// [5]int is Array []int is Slice is the value or type of this
	fmt.Printf("%T\n", x)
	fmt.Println(len(x)) // Print the length of the slice using len()

	fmt.Println("********one End*********")
	two()
	// fmt.Println("********two End*********")
	// three()
	// fmt.Println("********three End*********")
	// four()
	// fmt.Println("********four End*********")
}

func two() {
	x := []int{7, 2, 8, -9, 4, 0, 15, 22} // Slicing on a slice
	fmt.Println(x[:len(x)/2])
	fmt.Println(x[len(x)/2:])
}

func four() {
	x := []int{0, 1, 2, 3, 4}
	fmt.Println(x)
	x = append(x, 5, 6, 7)
	fmt.Println(x)

	y := []int{8, 9, 10}
	y = append(y, x[2:]...)
	fmt.Println(y)

	z := []int{11, 12, 13}
	z = append(x, z[:2]...) // This is interchangeable
	fmt.Println(z)
}

func three() {
	x := []int{42, 43, 44, 45, 46, 47, 48, 49, 50, 51}
	// [42, 43, 44, 48, 49, 50, 51]
	x = append(x[:3], x[6:]...)
	fmt.Println(x)
}
