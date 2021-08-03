package main

import "fmt"

type realMe struct {
	model string
}

func main() {
	m1 := realMe{
		model: "narzo30pro",
	}
	fmt.Println(m1)
	s := ChangeMe(&m1)
	fmt.Println(s)

}

// ChangeMe method takes the pointer values and change the underlying address value
func ChangeMe(p *realMe) string {
	p.model = "narzo20"
	s1 := p.model
	return s1
	// (*p).model = "realMe7"
}
