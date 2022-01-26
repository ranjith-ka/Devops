package main

import "fmt"

func main() {
	fmt.Println("Print the func retrun int:", foo())

	// h := bar()
	// fmt.Println(h())
	// i := scope()
	// fmt.Println(i)
	// fmt.Println(i)

	// // incrementor with function scope variable defined at the top

	// j := incrementor()
	// fmt.Println(j())
	// fmt.Println(j())
	// fmt.Println(j())
	// fmt.Println(j())
}

// Func return int
func foo() func() int {
	return func() int {
		return 43
	}

}

// func return  string
func bar() func() string {
	return func() string {
		return "hi"
	}
}

// Closer example, scope of the variable in the functions

func scope() int {
	var x int
	x++
	return x
}

// function return a func for a incrementor
func incrementor() func() int {
	var y int
	return func() int {
		y++
		return y
	}
}
