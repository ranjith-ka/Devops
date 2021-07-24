package main

import "fmt"

func main() {

	g := func() {
		for i := 0; i <= 10; i++ {
			fmt.Println(i)
		}
	}
	g()

	//Assign the function to a variable

	f := g
	f()
}
