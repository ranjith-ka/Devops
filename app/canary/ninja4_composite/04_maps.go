package main

import (
	"encoding/json"
	"fmt"
)

// Map introduction

func main() {

	six()
	seven()
}

func six() {

	m := map[string][]string{
		"Ramayana":    {"Ram", "Laxman", "Bharata", "Shatrughna"},
		"Mahabharata": {"Abhimanyu", "Balarama", "Draupadi", "Madri"},
	}
	fmt.Println(m)

	m["Silappatikaram"] = []string{"Kannaki", "Kovalan", "Madhavi"}

	fmt.Println(m)

	delete(m, "Mahabharata") // Delete in MAP

	fmt.Println(m)

	for k, v := range m {
		fmt.Println(k)
		for j, v2 := range v {
			fmt.Printf("\t%d \t  %s \n", j, v2)
		}
	}

	y := map[string]map[string]int32{
		"1gbdatablock": {
			"addons": 0,
			"newadd": 1,
		},
		"example": {
			"test": 1,
		},
	}

	z := map[string][]string{
		"1gbdatablock": {
			"addons",
			"test",
		},
		"example": {
			"test",
		},
	}
	fmt.Println(z)

	fmt.Println(y)
	fmt.Println(len(y))
	tm := len(m)
	for tm <= 2 {
		fmt.Println(y)
		tm++
	}
}

// https://www.youtube.com/watch?v=NDsbe1gNQG0&list=PLJ7-HiqskdZKn03J2X0y37agMrYCnhvEU&index=8

func seven() {

	Phone := make(map[string]int)

	m :=
		map[string]int{
			"Realme":  2021,
			"Oppo":    2022,
			"Samsung": 2021,
		}

		// Reassign the values
	Phone = m

	fmt.Println(Phone)

	fmt.Println("#### Cant expect which one comes first sometimes Oppo or Realme")
	for PhoneName, Year := range m {
		fmt.Printf("%s, %d\n", PhoneName, Year) // Cant expect which one comes first sometimes Oppo or Realme
	}

	//Verify the map with value
	if v, ok := Phone["Oppo"]; ok {
		fmt.Println("value:", v)

		delete(m, "Oppo")
	}

	// Convert one type of data structures to other
	data, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s", data)
}
