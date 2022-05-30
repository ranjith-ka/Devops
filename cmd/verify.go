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
		verfiyRandomJoke(url)
	},
}

func verfiyRandomJoke(Scopedurl string) string {
	resp, err := http.Head(Scopedurl)
	if err != nil {
		log.Println("Error", Scopedurl, err)
		return "404 NOK"
	}
	fmt.Println(resp.Status, Scopedurl)
	resp.Body.Close()
	return resp.Status
}
