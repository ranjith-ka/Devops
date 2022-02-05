package main

import (
	"fmt"
	"strings"
)

func main() {
	str1 := "Ooeno"
	str2 := "Qppo"

	s := strings.Compare(str1, str2)
	q := strings.Contains(str1, "is")

	fmt.Print(s, q)
	fmt.Println(strings.Index(str1, "o"))
	fmt.Println("Join Function: ", strings.Join([]string{"Australia", " ", "Japan", "Germany"}, ""))
}
