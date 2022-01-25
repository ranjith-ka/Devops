package main

import "fmt"

type Employee struct {
	Name string
}

// Interface to accept any structure and case to switch as per data.

func print(data interface{}) {
	switch val := data.(type) {

	case string:
		fmt.Println(val, "is a string")
	case int:
		fmt.Println(val, "int")
	case Employee: // Checking the struct via interface type casting(not a casting)
		fmt.Println(val.Name)
	}
}

func main() {
	x := "test"
	print(x)
	print(Employee{Name: "Keepcoding"})
}
