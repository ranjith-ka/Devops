package main

import (
	"encoding/json"
	"fmt"
)

type Phone struct {
	Model          string   `json:"model"`
	Cost           int      `json:"cost"`
	Brand          string   `json:"brand"`
	AvailablePlace []string `json:"available_places"`
}

func main() {
	p1 := Phone{
		Model: "narzo20pro",
		Cost:  200,
		Brand: "realme",
		AvailablePlace: []string{
			"bangalore",
			"karaikal",
		},
	}
	p2 := Phone{
		Model: "narzo30pro",
		Cost:  300,
		Brand: "realme",
		AvailablePlace: []string{
			"bangalore",
			"karaikal",
		},
	}
    p3 := Phone{
        Model: "Reno6",
        Cost: 500,
        Brand: "Oppo",
        AvailablePlace: []string{
            "Thrissur",
            "pattikad",
        },
    }

	mob := []Phone{p1, p2, p3}
	fmt.Println(mob)

	b, err := json.Marshal(mob)
	if err != nil {
		fmt.Println(err)
	}
	//Printing the byte slice
    fmt.Printf("%T\n", b)
	fmt.Println(b)
	// Just print as
	fmt.Println(string(b))
}
