package main

import (
"fmt"
    "math"
    "strings"
)

type MyStr string

type MyFloat float64

func (s MyStr) Uppercase() string {
    return strings.ToUpper(string(s))
}

func (f MyFloat) Abs() float64 {
    if f < 0 {
        return float64(-f)
    }
    return float64(f)
}

func main() {
    s := MyStr("test")  // Creating alias before using the methods.

    fmt.Println(s.Uppercase()) // To define methods on a type you don’t “own”, you need to define an alias for the type you want to extend:

    fmt.Println(MyStr("fun").Uppercase())

    f := MyFloat(-math.Sqrt2) // Convert to the specific type, then do the expression
    fmt.Println(f.Abs())
}
