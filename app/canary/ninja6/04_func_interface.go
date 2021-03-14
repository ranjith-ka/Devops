package main

import "fmt"

type moto struct {
	Brand string
	Model string
	price float64
	Tmp   string // not required for same data structures
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

// print the Interface type
func describe(i offer) {
	fmt.Printf("(%v, %T)\n", i, i)
}

func do(i interface{}) {
	switch v := i.(type) {
	case moto:
		fmt.Println("this is Moto", v.Model)
	case oppo:
		fmt.Println("this is oppo", v.Model)
	}
}
func main() {

	m1 := moto{
		Brand: "moto",
		Model: "v3",
		price: 20000,
		Tmp:   "ok",
	}
	m2 := oppo{
		Brand: "oppo",
		Model: "reno v3",
		price: 14000,
	}

	defer discount(m1) // defer to execute at last
	defer do(m1)
	defer describe(m1) // hold underlying type
	defer discount(m2)
	defer do(m2)
	defer describe(m2)

	// anonymous func
	func() {
		for i := 10; i <= 10; i++ {
			fmt.Println(i)
		}
		fmt.Println("done")
	}() // ending
}
