package main

import "fmt"

// https://www.golangprograms.com/go-language/functions.html

//An anonymous function is a function that was declared without any named identifier to refer to it. Anonymous functions can accept inputs and return outputs, just as standard functions do.

func main() {
	g := func() {
		for i := 0; i <= 10; i++ {
			fmt.Println(i)
		}
	}
	g()

	// Assign the function to a variable
	// f := g
	// f()

	SumPlusOne := func(i ...int) int {
		sum := 0
		for _, v := range i {
			sum += v
		}
		return sum + 1
	}

	fmt.Println("Lets print this:", SumPlusOne(2, 1, 5))

	l := 20
	b := 30

	// Closures are a special case of anonymous functions. Closures are anonymous functions which access the variables defined outside the body of the function.

	func() {
		var area int
		area = l * b
		fmt.Println(area)
	}()

}
