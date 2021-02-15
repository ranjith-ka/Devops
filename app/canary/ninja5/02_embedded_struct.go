// embedded structs

package main

import "fmt"

type vechile struct {
	doors int
	color string
}

type truck struct {
	vechile
	fourWheel bool
}

type sedan struct {
	vechile
	luxury bool
}

func main() {
	t := truck{
		vechile: vechile{
			doors: 2,
			color: "Grey",
		},
		fourWheel: true,
	}
	fmt.Println(t)

	s := sedan{
		vechile: vechile{ // embedding the structs in other structs.
			doors: 4,
			color: "black",
		},
		luxury: true,
	}
	fmt.Println(s)
}
