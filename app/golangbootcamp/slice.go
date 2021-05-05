package main

import "fmt"

func main(){
    s := make([]string, 3)
    s[0] = "test1"
    s[1] = "test2"
    s[2] = "test3"

    fmt.Printf("%q", s)
    fmt.Println(s)

    for _, v := range s{
        if v == "test1" {
            fmt.Println(v)
        }
    }

    // initialize with empty slice of int
    pow := make([]int, 10)
    for m := range pow{
        pow[m] = 1 << uint(m)
        if pow[m] > 200 {
            break
        }
    }
    fmt.Println(pow)
}
