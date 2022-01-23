// embedded structs

package main

import "fmt"

type vehicle struct {
	doors int
	color string
}

type truck struct {
    vehicle
	fourWheel bool
}

type sedan struct {
    vehicle
	luxury bool
}

func main() {
	t := truck{
		vehicle: vehicle{
			doors: 2,
			color: "Grey",
		},
		fourWheel: true,
	}
	fmt.Println(t)

	s := sedan{
		vehicle: vehicle{ // embedding the structs in other structs.
			doors: 4,
			color: "black",
		},
		luxury: true,
	}
	fmt.Println(s)
}
