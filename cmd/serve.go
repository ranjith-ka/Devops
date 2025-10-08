/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func hello(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "Here is my first http program")
	if err != nil {
		fmt.Printf("%v", err)
	}
}

func hello2(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "Here is my second http program guys")
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

func myjoke(w http.ResponseWriter, req *http.Request) {
    _, err := fmt.Fprint(w, getRandomJoke(req.Context()))
	if err != nil {
		fmt.Printf("%v", err)
	}
}

func myjoke2(w http.ResponseWriter, req *http.Request) {
    _, err := fmt.Fprint(w, getRandomJokeWithLLMStudio(req.Context()))
	if err != nil {
		fmt.Printf("%v", err)
	}
}

func serve() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Welcome to my website!")
		if err != nil {
			fmt.Printf("%v", err)
		}
	})
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/hello2", hello2)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/joke", myjoke)
	http.HandleFunc("/joke2", myjoke2)

	fmt.Println("Server up and running....")
	log.Print(http.ListenAndServe(":8080", nil))
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
