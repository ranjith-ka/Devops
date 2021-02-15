package main

import "fmt"

func main() {
	s := struct {
		first     string
		friends   map[string]int
		favDrinks []string
	}{
		first: "James",
		friends: map[string]int{
			"test1": 555,
			"test2": 777,
			"test3": 888,
		},
		favDrinks: []string{
			"ayaa",
			"oyaa",
		},
	}
	fmt.Println(s)
	fmt.Println(s.first)
	fmt.Println(s.friends)
	fmt.Println(s.favDrinks)

	for i, v := range s.friends {
		fmt.Println(i, v)
	}
	for j, v := range s.favDrinks {
		fmt.Println(j, v)
	}
}
