package main

import "fmt"

func main() {
	// Assign the values for new variables
	i, j := "gg", "yy"

	// Assign variables via pointer
	p := &i
	// Assign pointer variable to other variable
	y := &p

	// de-referencing the address to get the value
	fmt.Println(*p, j, i)
	fmt.Println(p)
	fmt.Println("Printing this")

	// de-reference the double pointer reference
    fmt.Println(**y)
	// change the value via pointer reference
	*p = "Change me"

	fmt.Println(i)
}
