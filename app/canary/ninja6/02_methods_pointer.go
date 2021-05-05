package main

import (
	"fmt"

	"github.com/dariubs/percent"
)

type Cost struct {
	X, Y float64
}

// FestivalOffer to add some offers during certain time.
// Pointer methods are used to change the values in that location instead working on the copy.
func (c *Cost) FestivalOffer(f float64) {
	c.X = c.X - percent.PercentFloat(f, c.X)
	c.Y = c.Y - percent.PercentFloat(f, c.Y)
}

func (c Cost) AddDiscount(f float64) (float64, float64) {
	return c.X - f, c.Y - f
}

func main() {
	var v Cost
	v = Cost{2000, 9000}
	v.FestivalOffer(20)
	fmt.Println(v.AddDiscount(500))

}
