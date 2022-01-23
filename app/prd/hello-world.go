package main

import (
	"fmt"
	"net/http"
	"os/user"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Println(user.Current())
        _, err := fmt.Fprintf(w, "Welcome to my prod website!")
        if err != nil {
            return
        }
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

    err := http.ListenAndServe(":8081", nil)
    if err != nil {
        return
    }
}
