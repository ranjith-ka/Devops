package main

// TDD

// function to print total
func sum(xi ...int) int {
	total := 0
	for _, v := range xi {
		total += v
	}
	return total
}

// function return a func for a incrementor
func incrementor() func() int {
	var y int
	return func() int {
		y++
		return y
	}
}

func main() {

}
