package main

import (
	"fmt"
	"net/http"
	"os/user"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Println(user.Current())
		fmt.Fprintf(w, "Welcome to my prod website!")
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", nil)
}
