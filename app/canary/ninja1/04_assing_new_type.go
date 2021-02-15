package main

import "fmt"

type mat int

var x mat
var y int

func main() {
	x = 12
	fmt.Println(x)
	fmt.Printf("%T\n", x)
	y = int(x)
	fmt.Println(y)
	fmt.Printf("%T\n", y)
}
