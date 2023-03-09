package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Create type to retrun the higher order function

// rune is an alias for int32 and is equivalent to int32 in all ways. It is used, by convention, to distinguish character values from integer values.
type RuneForRuneFunc func(rune) rune

func main() {

// Custom types unique to the enclosing package or Functions
type Count int
type StringMap map[string]string
type FloatChan chan float64



var i Count = 12
i--
fmt.Println(i)
sm := make(StringMap)
sm["key1"] = "Value1"
sm["key2"] = "Value2"
fmt.Println(sm)

// Channel
fc := make(FloatChan, 1)
fc <- 2.4534545345
fmt.Println(<-fc)

var removePunc RuneForRuneFunc

pharses := []string{"Day; :::is so long", "Get ;;;;;ready"}

fmt.Println(removePunc, pharses)

removePunc = func(char rune) rune{
	if unicode.Is(unicode.Terminal_Punctuation, char) {
		return -1
	}
	return char
}
processPharses(pharses, removePunc)
}

func processPharses(pharses []string, function RuneForRuneFunc){
	for _, phase := range pharses {
		fmt.Println(strings.Map(function, phase))
	}
}


