package main

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

	return wordCount
}

func main() {
	wc.Test(WordCount)
	fmt.Println(WordCount("test this test"))
}
