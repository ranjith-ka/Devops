// https://blog.golang.org/defer-panic-and-recover

// Check out the code and see the control flow here.
package main

import "fmt"

func main() {
    f()  // 1. call the function f to execute

    // 17. panic is recovered and print normally.
    fmt.Println("Returned normally from f.")
}

func f() {
    // 2. defer will execute at the end Last in first out defer statement.
    // Deferred function calls are executed in Last In First Out order after the surrounding function returns.

    // 15. defer func is waiting to execute , since its panic
    defer func() {
        // 16. a call to recover will capture the value given to panic and resume normal execution.
        if r := recover(); r != nil {
            fmt.Println("Recovered in f", r)
        }
    }()

    // 3. Just print here
    fmt.Println("Calling g.")
    // 4. Call the func g with argument 0
    g(0)

    // 18. Never return normally from g()
    fmt.Println("Returned normally from g.")
}

// 5. g starts the execution
func g(i int) {
    // 10. now g i s greater than 3,
    if i > 3 {
        // 11. Print below statement
        fmt.Println("Panicking!")
        // 12. Current go routine is panicking  now with g(4)
        panic(fmt.Sprintf("%v", i))
    }

    // 6. g less than 3, just print format
    // 7. Defer will print after this current execution.
    // 13. Print all the defer in the function g() with last in first out, so Defer g with 3, remember the panicking is still happening
    defer fmt.Println("Defer in g", i)
    // 14. complete the defer and go back to the calling function.

    // 8. Print the current value in g as 0,1,2,3
    fmt.Println("Printing in g", i)

    // 9. Increment g with 1 and call the function g() again
    g(i + 1)
}
