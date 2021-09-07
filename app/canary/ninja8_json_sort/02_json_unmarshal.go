package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net"
    "net/http"
    "time"
)



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

type idp struct {
    Realm           string `json:"realm"`
    PublicKey       string `json:"public_key"`
    TokenService    string `json:"token-service"`
    AccountService  string `json:"account-service"`
    TokensNotBefore int    `json:"tokens-not-before"`
}

// Adding the DailContext for Keycloak response and TLS timeout condition
var netTransport = &http.Transport{
    DialContext: (&net.Dialer{
        Timeout:   15 * time.Second,
        KeepAlive: 30 * time.Second,
    }).DialContext,
    TLSHandshakeTimeout: 5 * time.Second,
}
var netClient = &http.Client{
    Timeout: time.Second * 10,
    Transport: netTransport,
}

// getContent to copy the data from
func getContent(target *idp) error {
	url := `http://localhost:8080/auth/realms/Rule`
	res, err := netClient.Get(url)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
        err := Body.Close()
        if err != nil {

        }
    }(res.Body)
	return json.NewDecoder(res.Body).Decode(&target)
}

func main() {

    AllProduct := new(idp)
    err := getContent(AllProduct)
    if err != nil {
        return
    }

    fmt.Println(AllProduct.PublicKey)
	// Create a JSON structure type and convert the go type and print using the tags.

	// fmt.Println(AllProduct)
	// fmt.Println(AllProduct.Additionalfeatures)
}
