package main

import "fmt"

// func (r receiver) identifier(parameters) (return(s)) { code }

func main() {
	foo1()
	bar2("ooi", "gg")

	// x1 := rapper()
	// fmt.Println(x1)
	fmt.Println("rk")
	fmt.Println(rapper()) // Call function directory without variable to assign

	// 2 values in return
	// x2, x3 := test()
	fmt.Println(test())

}

// foo Just a function to call
func foo1() {
	fmt.Println("Hey this is foo")
}

// String Return

func bar2(s, e string) {
	fmt.Println("bar", s, e)
}

func rapper2() string {
	return "hi"
}

// Int return
func rapper() int {
	return 42
}

// Return String + int
func test() (int, string) {
	return 42, "testing this"
}
