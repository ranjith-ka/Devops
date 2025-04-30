package wc

import "fmt"

func main() {
	i, j := "gg", "yy"
	p := &i
	fmt.Println(*p, j, i)
	fmt.Println(p)
	// dereference the pointer location
	*p = "changeme"
	fmt.Println(i)
}
