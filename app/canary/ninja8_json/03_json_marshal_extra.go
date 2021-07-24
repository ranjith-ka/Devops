package main

import (
    "encoding/json"
    "fmt"
)

//easyOne struct type is known to us the values we are going to parse
// Note: Some and Values should be in CAPS to export
type easyOne struct{
    Some string
    Value string
}

func main(){
    myJsonString := `{"some": "super", "value": "got-ehm"}`
    var v easyOne
    err := json.Unmarshal([]byte(myJsonString), &v)
    if err != nil {
        return
    }
    fmt.Printf("%s\n%s\n", v.Some,v.Value)

    //Unknown struct if we have in JSON
    birdJson := `{"birds":{"pigeon":"likes to perch on rocks","eagle":"bird of prey"},"animals":"none"}`
    var result map[string]interface{}
    err = json.Unmarshal([]byte(birdJson), &result)
    if err != nil {
        return
    }
    fmt.Println(result)
    // The object stored in the "birds" key is also stored as
    // a map[string]interface{} type, and its type is asserted from
    // the interface{} type
    birds := result["birds"].(map[string]interface{})

    for key, value := range birds {
        // Each value is an interface{} type, that is type asserted as a string
        fmt.Println(key, value.(string))
    }

}
