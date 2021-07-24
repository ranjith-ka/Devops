package main

import "fmt"

type mobile struct {
	Brand    string
	colour   string
	Model    string
	location []string
}

func main() {
	m1 := mobile{
		Brand:  "Oppo",
		colour: "Black",
		Model:  "Reno5 Pro",
		location: []string{
			"India",
			"China",
		},
	}

	m2 := mobile{
		Brand:  "Honor",
		colour: "Grey",
		Model:  "9 Lite",
		location: []string{
			"India",
			"China",
			"US",
		},
	}
	fmt.Println(m1.Brand)
	fmt.Println(m1.Model)
	for i, v := range m1.location {
		fmt.Println(i, v)
	}

	fmt.Println(m2.Brand)
	fmt.Println(m2.Model)
	for i, v := range m2.location {
		fmt.Println(i, v)
	}

	fmt.Println("*********End*********")

	// Map for other type in Struct

	p := map[string]mobile{
		m1.Brand: m1,
		m2.Brand: m2,
	}
	fmt.Println(p)

	// Range over location
	for _, v := range p {
		fmt.Println(v.Brand)
		for j, v2 := range v.location {
			fmt.Println(j, v2)
		}
	}
}
