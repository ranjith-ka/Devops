package main

import "fmt"

type vechile struct {
	name   string
	color  string
	lights rune
}

func (v vechile) DisplayLight() {
	fmt.Printf("The no of light in %s is: %d\n", v.name, v.lights)
}

func main() {
	v1 := vechile{
		name:   "i20",
		color:  "black",
		lights: 5,
	}

	v2 := vechile{
		name:   "i60",
		color:  "black",
		lights: 4,
	}

	v1.DisplayLight()
	v2.DisplayLight()
}
