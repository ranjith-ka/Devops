package main

import "fmt"

type moto struct {
	Brand string
	Model string
	price float64
}

type oppo struct {
	Brand string
	Model string
	price float64
}

// func (r receiver) identifier(parameters) (return(s)) { code }

func (m moto) curentOffer() float64 {
	return m.price - (m.price * 0.25)
}

func (o oppo) curentOffer() float64 {
	return o.price - (o.price * 0.15)
}

type offer interface {
	curentOffer() float64
}

func discount(o offer) {
	fmt.Println(o.curentOffer())
}

func main() {

	m1 := moto{
		Brand: "moto",
		Model: "v3",
		price: 20000,
	}
	m2 := oppo{
		Brand: "oppo",
		Model: "reno v3",
		price: 14000,
	}

	defer discount(m1) // defer to execute at last
	defer discount(m2)

	// anonymous func

	func() {
		for i := 10; i <= 10; i++ {
			fmt.Println(i)
		}
		fmt.Println("done")
	}() // ending
}
