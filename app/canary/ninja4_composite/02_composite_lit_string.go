package main

import "fmt"

func main() {

	x := make([]string, 1, 36) // length and capacity min 1 and max 36 , this use case
	x = []string{`Andhra Pradesh`, `Arunachal Pradesh`, `Assam`, `Bihar`, `Chhattisgarh`, `Goa`, `Gujarat`, `Haryana`, `Himachal Pradesh`, `Jammu and Kashmir`, `Jharkhand`, `Karnataka`, `Kerala`, `Madhya Pradesh`, `Maharashtra`, `Manipur`, `Meghalaya`, `Mizoram`, `Nagaland`, `Odisha`, `Punjab`, `Rajasthan`, `Sikkim`, `Tamil Nadu`, `Telangana`, `Tripura`, `Uttar Pradesh`, `Uttarakhand`, `West Bengal`, `Andaman and Nicobar Islands`, `Chandigarh`, `Dadra and Nagar Haveli`, `Daman and Diu`, `National Capital Territory of Delhi`, `Lakshadweep`, `Pondicherry`}
	fmt.Println(x)
	fmt.Println(len(x))

	fmt.Println(cap(x))

	for i := 0; i < len(x); i++ {
		fmt.Println(i, x[i])
	}
	fmt.Println("****End of First Program****")

	five()
}

// Second program for range with String  composite literals
func five() {
	xs1 := []string{"one", "two", "three"}
	xs2 := []string{"three", "four", "five", "six"}
	fmt.Println(xs1)
	fmt.Println(xs2)

	// Slice of slice if String

	xxs := [][]string{xs1, xs2}
	fmt.Println("Lets have fun")
	fmt.Println(xxs)

	for i, xs := range xxs {
		fmt.Println("Record:", i)
		for j, val := range xs {
			fmt.Printf("Index Position: %d \t value: \t %s \n", j, val)
		}
	}
}
