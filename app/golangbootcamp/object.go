package main

import "fmt"

type gadgets uint8

const (
    Camera gadgets = iota
    Bluetooth
    Media
    Storage
    VideoCalling
    Multitasking
    Messaging
)
type mobile struct {
    make string
    model string
}

type smartphone struct {
    gadgets gadgets
}

func (s *smartphone) launch() {
    fmt.Println ("New Smartphone Launched:", "with total gadgets", s.gadgets )
}

type android struct {
    mobile
    smartphone
    waterproof string
}
func (a *android) samsung() {
    fmt.Printf("%s %s\n",
        a.make, a.model)
}

type iphone struct {
    mobile
    smartphone
    sensor int
}

func (i *iphone) apple() {
    fmt.Printf("%s %s %d\n",
        i.make, i.model, i.gadgets)
}


func main() {
    t := &iphone {}
    t.make ="Samsung"
    t.model ="Galaxy J7 Prime"
    t.gadgets = Camera+Bluetooth+Media+Storage+VideoCalling+Multitasking+Messaging
    t.launch()
    t.apple()

    m2 := &mobile{
        make: "apple",
        model: "SE"}

    f := &android{*m2,
        smartphone{
            Camera+Bluetooth+Messaging,
        },
    "test",
    }

    fmt.Println(f.model)
    fmt.Println(m2)
    f.launch()
    f.samsung()

}
