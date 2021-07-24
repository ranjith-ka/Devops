package wc

import (
"fmt"
"time"
)

func main() {
    fmt.Println("When's Saturday?")
    today := time.Now().Weekday()
    //timeNow := time.Unix(2,1)
    fmt.Printf("today:%s\n", today)
    switch time.Thursday {
    case today + 0:
        fmt.Println("Today.")
    case today + 1:
        fmt.Println("Tomorrow.")
    case today + 2:
        fmt.Println("In two days.")
    default:
        fmt.Println("Too far away.")
    }
}
