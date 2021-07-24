package main

import "fmt"

func main() {
	a := 42
	fmt.Println(a)
	fmt.Printf("%T\n", a)  // Type of variable
	fmt.Printf("%b\n", a)  // Binary
	fmt.Printf("%d\n", a)  // decimal
	fmt.Printf("%#x\n", a) // Hexadecimal
}
