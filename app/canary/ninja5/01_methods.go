package main

import "fmt"

type Brand struct {
	X, Y string
}

// Phone to return if brand is Moto
func (b Brand) Phone() string {
	if b.X == "moto" {
		return b.Y
	}
	return b.X
}

func main() {
	var v Brand // To ignore the error in vscode, but this is the defined type
	v = Brand{"moto", "return me if moto in first field"}
	fmt.Println(v.Phone())
}
