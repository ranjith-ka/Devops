package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	n, err := fmt.Fprint(w, "Here is my first http program")
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Println(n)
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			n, err := fmt.Fprintf(w, "%v: %v\n", name, h)
			if err != nil {
				fmt.Printf("%v", err)
			}
			fmt.Println(n)
		}
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		n, err := fmt.Fprintf(w, "Welcome to my canary website!")
		if err != nil {
			fmt.Printf("%v", err)
		}
		fmt.Println(n)
	})
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	log.Print(http.ListenAndServe(":8080", nil))
}
