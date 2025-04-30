package main

import "fmt"

type mat int

func main() {
	var x mat
	var y int
	x = 12
	fmt.Println(x)
	fmt.Printf("%T\n", x)
	y = int(x)
	fmt.Println(y)
	fmt.Printf("%T\n", y)
}
