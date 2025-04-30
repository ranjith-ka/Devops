package main

import (
	"fmt"

	t "github.com/dariubs/percent"
)

type Mymobile struct {
	Brand string
	model []string
	yom   int
	price int
}

func (m *Mymobile) motoOffer() int {
	if m.Brand == "moto" {
		m.price = m.price - 200
	}
	return m.price
}

func (m Mymobile) oppoOffer() int {
	if m.Brand == "oppo" {
		for _, v := range m.model {
			if v == "reno" {
				m.price = m.price - int(t.Percent(20, m.price))
			}
		}
	}
	return m.price
}

func main() {
	s := Mymobile{
		Brand: "moto",
		model: []string{
			"model1",
			"model2",
		},
		yom:   2008,
		price: 12000,
	}
	fmt.Printf("After offer for moto %d\n", s.motoOffer())
	fmt.Println(s.price)

	p := Mymobile{
		Brand: "oppo",
		model: []string{
			"reno",
			"model2",
		},
		yom:   2008,
		price: 12000,
	}
	fmt.Printf("After offer for oppo %d", p.oppoOffer())

}
