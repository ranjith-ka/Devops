package main

import (
    "fmt"
)

func main() {
	x := []int{1, 5, 3, 4, 2} // Slice

	// x := [5]int{0, 1, 4, 3, 2}  //  Array
	fmt.Println(x)
	for i, v := range x {
		fmt.Println(i, v)
	}
	// [5]int is Array []int is Slice is the value or type of this
	fmt.Printf("%T\n", x)

	fmt.Println("********one End*********")
	two()
	fmt.Println("********two End*********")
	four()
	fmt.Println("********three End*********")
	three()
	fmt.Println("********four End*********")
}

func two() {
	x := []int{1, 2, 3, 4, 5, 6, 7, 8, 9} // Slicing on a slice
	fmt.Println(x)
	fmt.Println(x[:5])
	fmt.Println(x[5:])
	fmt.Println(x[2:7])
}

func four() {
	x := []int{0, 1, 2, 3, 4}
	fmt.Println(x)
	x = append(x, 5, 6, 7)
	fmt.Println(x)

	y := []int{8, 9, 10}
	y = append(y, x[:2]...)
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
