package wc

import (
	"fmt"
	"strings"

	"golang.org/x/tour/wc"
)

func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	fmt.Println(words)
	wordCount := make(map[string]int)
	fmt.Println(wordCount)
	for i := range words {
		wordCount[words[i]]++
		fmt.Println(i)
	}
	wordCount["ok"]=21
    return wordCount
}

func main() {
	wc.Test(WordCount)
	fmt.Println(WordCount("test this test ok ok ok ok ok "))
}

func test() {
    var employee = map[string]int{}
    fmt.Println(employee)        // map[]
    fmt.Printf("%T\n", employee) // map[string]int
}
