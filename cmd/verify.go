/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "To verify the source URL for random joke",
	Long:  `Just try to verify the URL to check the source or knowloedge databse, so make sure we run something`,
	Run: func(cmd *cobra.Command, args []string) {
		verfiyRandomJoke()
	},
}

func verfiyRandomJoke() {
	resp, err := http.Head(url)
	if err != nil {
		log.Println("Error", url, err)
	}
	fmt.Println(resp.Status, url)
	resp.Body.Close()
}
