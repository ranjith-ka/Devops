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

// method set for above types
func (m *moto) priceChange() {
   m.price = m.price - 500
}

func (o *oppo) priceChange()  {
    o.price = o.price - 300
}

// func (r receiver) identifier(parameters) (return(s)) { code }

func (m moto) currentOffer() float64 {
    return m.price - (m.price * 0.25)
}

func (o oppo) currentOffer() float64 {
    return o.price - (o.price * 0.15)
}

type offer interface {
    currentOffer() float64
}

func discount(o offer) {
    fmt.Println(o.currentOffer())
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
    }
    m2 := oppo{
       Brand: "oppo",
       Model: "reno v3",
       price: 14000,
    }

     discount(m1) // defer to execute at last
     do(m1)
     describe(m1) // hold underlying type
     //discount(m2)
     //do(m2)
     //describe(m2)
     m1.priceChange()
     fmt.Println(m1)
     discount(m1)
     m2.priceChange()
    fmt.Println(m2)
}
