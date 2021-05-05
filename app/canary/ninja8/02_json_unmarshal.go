package main

import (
	"encoding/json"
	"fmt"
    "io"
    "net/http"
)

type product interface{} // use Interface to assume the unknown data to the structure.

// type product struct {
// 	Additionalfeatures string `json:"additionalFeatures"`
// 	Os                 string `json:"os"`
// 	Battery            struct {
// 		Type        string `json:"type"`
// 		Standbytime string `json:"standbytime"`
// 	} `json:"battery"`
// 	Camera struct {
// 		Features []string `json:"features"`
// 		Primary  string   `json:"primary"`
// 	} `json:"camera"`
// 	Connectivity struct {
// 		Bluetooth string `json:"bluetooth"`
// 		Cell      string `json:"cell"`
// 		Cell      string `json:"cell"`
// 		Gps       bool   `json:"gps"`
// 		Infrared  bool   `json:"infrared"`
// 		Wifi      string `json:"wifi"`
// 	} `json:"connectivity"`
// 	Description string `json:"description"`
// 	Display     struct {
// 		Screenresolution string `json:"screenResolution"`
// 		Screensize       string `json:"screenSize"`
// 	} `json:"display"`
// 	Hardware struct {
// 		Accelerometer    bool   `json:"accelerometer"`
// 		Audiojack        string `json:"audioJack"`
// 		CPU              string `json:"cpu"`
// 		Fmradio          bool   `json:"fmRadio"`
// 		Physicalkeyboard bool   `json:"physicalKeyboard"`
// 		Usb              string `json:"usb"`
// 	} `json:"hardware"`
// 	ID            string   `json:"id"`
// 	Images        []string `json:"images"`
// 	Name          string   `json:"name"`
// 	Sizeandweight struct {
// 		Dimensions []string `json:"dimensions"`
// 		Weight     string   `json:"weight"`
// 	} `json:"sizeAndWeight"`
// 	Storage struct {
// 		Hdd string `json:"hdd"`
// 		RAM string `json:"ram"`
// 	} `json:"storage"`
// }

// getContent to copy the data from
func getContent(target interface{}) error {
	url := `https://gist.githubusercontent.com/vivekjuneja/5369525/raw/88a6291fa9aa700868724751b4e7485c073e0289/product.json`
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
        err := Body.Close()
        if err != nil {

        }
    }(res.Body)
	return json.NewDecoder(res.Body).Decode(target)
}

func main() {

	AllProduct := new(product)
    err := getContent(AllProduct)
    if err != nil {
        return
    }
	fmt.Println(*AllProduct)

	// Create a JSON structure type and convert the go type and print using the tags.

	// fmt.Println(AllProduct)
	// fmt.Println(AllProduct.Additionalfeatures)
}
