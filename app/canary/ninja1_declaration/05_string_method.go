package main

import (
	"fmt"
	"strings"
)

func main() {
	str1 := "Ooeno"
	str2 := "Qppo"
	str3 := "Ooeno"

	s := strings.Compare(str1, str2)
	z := strings.Compare(str1, str3)
	q := strings.Contains(str1, "is")
	

	fmt.Println(s, z, q)
	fmt.Println(strings.Index(str1, "o"))
	fmt.Println("Join Function: ", strings.Join([]string{"Australia", " ", "Japan", "Germany"}, ""))
}
