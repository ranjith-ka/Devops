package main

import (
    "fmt"
    "sort"
)

type sortUser struct {
    First   string
    Last    string
    Age     int
    Sayings []string
}

type ByAge []sortUser

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

//TODO Other day add String to sort in this type

func main() {
    xi := []int{5, 8, 2, 43, 17, 987, 14, 12, 21, 1, 4, 2, 3, 93, 13}
    xs := []string{"random", "rainbow", "delights", "in", "torpedo", "summers", "under", "gallantry", "fragmented", "moons", "across", "magenta"}
    // Sort by type Int
    sort.Ints(xi)
    // Sort by type String
    sort.Strings(xs)
    fmt.Println(xi,"\n",xs)

    u1 := sortUser{
        First: "James",
        Last:  "Bond",
        Age:   12,
        Sayings: []string{
            "Shaken, not stirred",
            "Youth is no guarantee of innovation",
            "In his majesty's royal service",
        },
    }

    u2 := sortUser{
        First: "Miss",
        Last:  "Moneypenny",
        Age:   27,
        Sayings: []string{
            "James, it is soo good to see you",
            "Would you like me to take care of that for you, James?",
            "I would really prefer to be a secret agent myself.",
        },
    }

    u3 := sortUser{
        First: "M",
        Last:  "Hmmmm",
        Age:   54,
        Sayings: []string{
            "Oh, James. You didn't.",
            "Dear God, what has James done now?",
            "Can someone please tell me where James Bond is?",
        },
    }

    users := []sortUser{u1, u2, u3}
    fmt.Println(users)

    // sort.Sort required a data from interface, ByAge is a struct which implements the methods with signature
    // Default sort.Interface also implements this signature, hence sort.Sort with consider ByAge as interface since it
    // has same method signature. So we need to have all method signature to make sure we implement the interface.
    sort.Sort(ByAge(users))
    fmt.Println(users)
}
