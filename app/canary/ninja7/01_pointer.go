package main

import "fmt"

func main() {
	// Assign the values for new variables
	i, j := "gg", "yy"

	// Assign variables via pointer
	p := &i

	// de-referencing the address to get the value
	fmt.Println(*p, j, i)
	fmt.Println(p)

	// change the value via pointer reference
	*p = "changeme"

	fmt.Println(i)
}
