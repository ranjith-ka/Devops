package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func hello(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "Here is my first http program")
	if err != nil {
		fmt.Printf("%v", err)
	}
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			_, err := fmt.Fprintf(w, "%v: %v\n", name, h)
			if err != nil {
				fmt.Printf("%v", err)
			}
		}
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Welcome to my website!")
		if err != nil {
			fmt.Printf("%v", err)
		}
	})
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.Handle("/metrics", promhttp.Handler())
	// mux := http.NewServeMux()
	fmt.Println("Server up and running....")
	// http.ListenAndServe(":8080", apmhttp.Wrap(mux))
	log.Print(http.ListenAndServe(":8080", nil))
}
