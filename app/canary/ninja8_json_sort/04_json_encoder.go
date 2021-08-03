package main

import (
    "encoding/json"
    "fmt"
    "os"
)

type user struct {
    First   string
    Last    string
    Age     int
    Sayings []string
}

func main() {
    u1 := user{
        First: "James",
        Last:  "Bond",
        Age:   32,
        Sayings: []string{
            "Shaken, not stirred",
            "Youth is no guarantee of innovation",
            "In his majesty's royal service",
        },
    }

    u2 := user{
        First: "Miss",
        Last:  "Moneypenny",
        Age:   27,
        Sayings: []string{
            "James, it is soo good to see you",
            "Would you like me to take care of that for you, James?",
            "I would really prefer to be a secret agent myself.",
        },
    }

    u3 := user{
        First: "M",
        Last:  "Hmmmm",
        Age:   54,
        Sayings: []string{
            "Oh, James. You didn't.",
            "Dear God, what has James done now?",
            "Can someone please tell me where James Bond is?",
        },
    }

    users := []user{u1, u2, u3}

    // This Prints the list of string not in the JSON format.
    fmt.Println(users)

    // your code goes here
    // NewEncoder method for stream of data, it can print in the other interface
    // where in Marshal, need to be printed in strings.
    err := json.NewEncoder(os.Stdout).Encode(users)
    // Again os.Stdout have a same return type, hence we can the Type File in OS package to use as parameters.
    if err != nil {
        fmt.Println("We did something wrong and here's the error:", err)
    }

    b, err := json.Marshal(users)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(string(b))
}
