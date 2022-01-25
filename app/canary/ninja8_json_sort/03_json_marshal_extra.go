package main

import (
	"encoding/json"
	"fmt"
)

//easyOne struct type is known to us the values we are going to parse
// Note: Some and Values should be in CAPS to export
type easyOne struct {
	Some  string
	Value string
}

type mojo []struct {
	First   string   `json:"First"`
	Last    string   `json:"Last"`
	Age     int      `json:"Age"`
	Sayings []string `json:"Sayings"`
}

func main() {
	myJsonString := `{"some": "super", "value": "got-ehm"}`
	var v easyOne
	err := json.Unmarshal([]byte(myJsonString), &v)
	if err != nil {
		return
	}
	fmt.Printf("%s\n%s\n", v.Some, v.Value)

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

	// Ninja 8 implementation
	s := `[{"First":"James","Last":"Bond","Age":32,"Sayings":["Shaken, not stirred","Youth is no guarantee of innovation","In his majesty's royal service"]},{"First":"Miss","Last":"Moneypenny","Age":27,"Sayings":["James, it is soo good to see you","Would you like me to take care of that for you, James?","I would really prefer to be a secret agent myself."]},{"First":"M","Last":"Hmmmm","Age":54,"Sayings":["Oh, James. You didn't.","Dear God, what has James done now?","Can someone please tell me where James Bond is?"]}]`
	var b mojo
	err = json.Unmarshal([]byte(s), &b)
	if err != nil {
		return
	}
	fmt.Println(b)
	for i, person := range b {
		fmt.Println("Person #", i)
		fmt.Println("\t", person.First, person.Last, person.Age)
		for _, saying := range person.Sayings {
			fmt.Println("\t\t", saying)
		}
	}

}
