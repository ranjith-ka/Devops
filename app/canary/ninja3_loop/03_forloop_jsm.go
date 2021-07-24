package main

import (
	"fmt"
)

func main() {
	for i := 0; i <= 20; i++ {
		fmt.Println(i)
	}

	fmt.Println("Next Program2")
	two()
	fmt.Println("Next Program 3")
	three()
	fmt.Println("Next Program 4")
	four()
	fmt.Println("Next Program 5")
	five()
	//fmt.Println("Next Program 6")
	//six()
}

func two() { //a tab with \t
	for i := 65; i <= 90; i++ { //a new line with \n
		fmt.Println(i)           //hex value with %#x
		for j := 0; j < 1; j++ { //decimal value with %d
			fmt.Printf("\t%#U\n", i) //unicode code point with %#U
		}
	}
}

func three() {
	bd := 1995
	for bd <= 2021 {
		fmt.Println(bd) // for condition()
		bd++
	}
}

func four() {
	bd := 1995
	for {
		if bd > 2020 {
			break
		} // for loop with break statememt....
		fmt.Println(bd)
		bd++
	}
}

func five() {
	for i := 10; i <= 100; i++ { //%v	the value in a default formatwhen printing structs, the plus flag (%+v) adds field names
		fmt.Printf("When %v is divided by 4, the remainder is %v\n", i, i%4)
	}
}
