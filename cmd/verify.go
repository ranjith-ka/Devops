// Package cmd /*

package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ranjith-ka/Devops/database"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "To verify the source URL for random joke",
	Long:  `Just try to verify the URL to check the source or knowloedge databse, so make sure we run something`,
	Run: func(cmd *cobra.Command, args []string) {
		verifyRandomJoke(url)
	},
}

func verifyRandomJoke(ScopedUrl string) string {

	resp, err := http.Head(ScopedUrl)
	if err != nil {
		log.Println("Error", ScopedUrl, err)
		return "404 NOK"
	}
	fmt.Println(resp.Status, ScopedUrl)
	err = resp.Body.Close()
	if err != nil {
		return ""
	}

	database.SetupDB()

	return resp.Status
}
