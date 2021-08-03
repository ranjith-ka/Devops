package main

import (
    "fmt"
    "sync"
)

func main(){
    var wg sync.WaitGroup
    wg.Add(2)

    go func(){
        fmt.Println("Realme offers here")
        wg.Done()
    }()

    go func(){
        fmt.Println("Oppo offers here")
        wg.Done()
    }()

    // Uncomment the Wait to check the new offers, else you might miss the offers
    wg.Wait()
    fmt.Println("Check the new mobiles with offers, make sure you didn't miss in concurrency")
}
