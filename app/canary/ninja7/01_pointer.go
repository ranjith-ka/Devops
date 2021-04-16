package main

import "fmt"

func main() {
	i, j := "gg", "yy"
	p := &i
	fmt.Println(*p, j, i)
	fmt.Println(p)
	*p = "changeme"
	fmt.Println(i)
}
