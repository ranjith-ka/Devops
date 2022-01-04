package main

import "fmt"

func main() {
	a := 42
	fmt.Println(a)
	fmt.Printf("%T\n", a)                       // Type of variable
	fmt.Printf("%b\n", a)                       // Binary
	fmt.Printf("%d\n", a)                       // decimal
	fmt.Printf("%#x\n", a)                      // Hexadecimal
	fmt.Printf("%T\t%b\t%d\t%#x\n", a, a, a, a) // Type of variable
	// Shifting 2^3 + 2^2 + 2^1 (left shift to add and right shift to reduce it)
	b := a >> 1
	fmt.Printf("%T\t%b\t%d\t%#x\n", b, b, b, b)
}
