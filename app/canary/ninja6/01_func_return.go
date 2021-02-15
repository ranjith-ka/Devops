package main

import "fmt"

// func (r receiver) identifier(parameters) (return(s)) { code }

func main() {
	foo()
	bar("ooi")

	// x1 := rapper()
	// fmt.Println(x1)
	fmt.Println("rk")
	fmt.Println(rapper()) // Call function direcoty without variable to assign

	// 2 values in return
	x2, x3 := test()
	fmt.Println(x2, x3)

}

func foo() {
	fmt.Println("Hey this is foo")
}

// String Return

func bar(s string) {
	fmt.Println("bar", s)
}

func rapper2() string {
	return "hi"
}

// Int retrun
func rapper() int {
	return 42
}

// Retrun String + int
func test() (int, string) {
	return 42, "testing this"
}
