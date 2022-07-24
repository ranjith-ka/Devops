package main

import (
	"fmt"

	t "github.com/dariubs/percent"
)

type Cost struct {
	X, Y float64
}

// FestivalOffer to add some offers during certain time.
// Pointer methods are used to change the values in that location instead working on the copy.
func (c *Cost) FestivalOffer(f float64) (float64, float64) {
	c.X = c.X - t.PercentFloat(f, c.X)
	c.Y = c.Y - t.PercentFloat(f, c.Y)

	return c.X, c.Y
}

func (c *Cost) AddDiscount(f float64) (float64, float64) {
	return c.X - f, c.Y - f
}

func main() {
	var v Cost
	v = Cost{2000, 9000}
	// Run the festival offer first
	fmt.Println(v.FestivalOffer(10))
	// And then add additional discont on top the discount
	fmt.Println(v.AddDiscount(500))

}
