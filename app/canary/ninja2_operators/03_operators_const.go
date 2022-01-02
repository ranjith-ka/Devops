package main

import (
	"bytes"
	"fmt"
)

const (
	aa     = 42
	ab int = 43
)

var buffer [256]byte

type sliceHeader struct {
	Length        int
	ZerothElement *byte
}

func main() {
	a := (42 == 42)
	b := (42 <= 43)
	c := (42 >= 43)
	d := (42 != 43)
	e := (42 < 43)
	f := (42 > 43)

	fmt.Println(a, b, c, d, e, f)
	fmt.Println(aa, ab)

	// just try some random examples for slices, https://go.dev/blog/slices
	slice2 := sliceHeader{
		Length:        50,
		ZerothElement: &buffer[100],
	}

	slice := buffer[100:150]

	// Just Print the slice for reference
	fmt.Println(slice)

	for range slice {
		if len(slice) == 0 {
			break
		}
		slice = slice[1 : len(slice)-1]
		fmt.Println(slice)
	}

	fmt.Println(slice2)
	slashPos := bytes.IndexRune(slice, '/')
	fmt.Println(slashPos)
}
