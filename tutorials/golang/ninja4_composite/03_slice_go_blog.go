package main

import (
	"bytes"
	"fmt"
)

var buffer [256]byte

type sliceHeader struct {
	Length        int
	ZerothElement *byte
}

func main() {

	// just try some random examples for slices, https://go.dev/blog/slices
	slice2 := sliceHeader{
		Length:        5,
		ZerothElement: &buffer[0],
	}

	slice := buffer[100:120]

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

	slice = buffer[100:120]

	fmt.Println(slice)
	for i := 0; i < len(slice); i++ {
		slice[i] = byte(i)
	}
	fmt.Println("before", slice)
	AddOneToEachElement(slice)
	fmt.Println("after", slice)

}

// passing slice to function
func AddOneToEachElement(slice []byte) {
	for i := range slice {
		slice[i]++
	}
}
