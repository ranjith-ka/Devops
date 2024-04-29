package main

import "fmt"

func main() {
	one()
	//two()
}

type Brand struct {
	X, Y string
}

// Phone to return if brand is Moto
func (b Brand) Phone() string {
	if b.X == "moto" {
		return b.X
	}
	return b.Y
}

func one() {

	var v Brand // To ignore the error in vscode, but this is the defined type
	v = Brand{"moto", "return me if moto in first field"}

	s := Brand{X: "ffg", Y: "test"}
	fmt.Println(v.Phone())
	fmt.Println(s.Phone())
}

// https://www.youtube.com/watch?v=HDF5Ol8jto0&list=PLJ7-HiqskdZKn03J2X0y37agMrYCnhvEU&index=11

type Celsius float32

type Fahrenheit float32

func (c Celsius) ToFahrenheit() Fahrenheit {
	return Fahrenheit(9.0/5.0*c + 32)
}

func two() {

	// conversion of float with celsius and call the method.
	f := Celsius(36.7).ToFahrenheit()
	fmt.Println(f)
}
