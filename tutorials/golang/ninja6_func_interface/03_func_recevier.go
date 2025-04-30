package main

import "fmt"

type mobile struct {
	Brand    string
	Model    string
	discount int
}

// func (r receiver) identifier(parameters) (return(s)) { code }

func (m *mobile) disc() {
	fmt.Println("Discount for the mobile", m.Brand, "is", m.discount, "%")
}

func main() {
	m1 := &mobile{
		Brand:    "Oppo",
		Model:    "Reno5 pro",
		discount: 20,
	}
	m2 := &mobile{
		Brand:    "Honor",
		Model:    "9 Lite",
		discount: 30,
	}
	m1.disc()
	m2.disc()
}
