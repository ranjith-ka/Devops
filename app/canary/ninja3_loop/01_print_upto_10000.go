package main

import "fmt"

func main() {
	for i := 0; i <= 10; i++ {
		fmt.Println(i)
	}
	for {
		fmt.Println("test this")
		break
	}
	s := 1
	for s <= 1 {
		fmt.Println(s)
		break
	}
}
