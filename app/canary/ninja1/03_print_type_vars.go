package main

import "fmt"

func main() {
    const f = "%T(%v)\n"
	x := 42
	y := "James Bond"
	z := true

	s := fmt.Sprintf("%d\t%s\t%t", x, y, z)
	r := fmt.Sprintf("%v score is %v", y, x)
	fmt.Println(s)
	fmt.Println(r)
    fmt.Printf(f, x, x)
}
