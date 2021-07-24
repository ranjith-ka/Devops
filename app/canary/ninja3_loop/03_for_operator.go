package main

import "fmt"

// Default operators
func main() {
	bd := 1975
	for bd <= 1985 {
		fmt.Println(bd)
		bd++
	}
	fmt.Println("Next Program2")
	two()
	fmt.Println("Next Program 3")
	three()
	fmt.Println("Next Program 4")
	four()
	fmt.Println("Next Program 5")
	five()
	fmt.Println("Next Program 6")
	six()
}

// If operator with for loop

func two() {
	bd := 2000
	for {
		if bd > 2017 {
			break
		}
		fmt.Println(bd)
		bd++
	}
}

// Modulus Operator

func three() {
	for i := 10; i <= 100; i++ {
		fmt.Printf("when %d is divded by 4, the remainder is %d\n", i, i%4)
	}
}

// If statement

func four() {
	x := "Ranjith KA"
	if x == "Ranjith KA" {
		fmt.Println(x)
	}
}

// If else if, else statement example

func five() {
	x := 27
	if x == 25 {
		fmt.Printf("Yes its %d\n", x)
	} else if x == 26 {
		fmt.Println(x)
	} else {
		fmt.Println("not 25 nor 26")
	}
}

// Switch statement

func six() {
	switch {
	case true:
		fmt.Println("print")
	case false:
		fmt.Println("skip")
	}
}
