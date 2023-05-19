/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var wellknown_configuration map[string]any

func fetch(cmd *cobra.Command, args []string) {
	var hostname string
	if len(args) == 0 {
		hostname = "http://127.0.0.1"
	} else {
		hostname = args[0]
	}

	well_known_endpoint := hostname + "/.well-known/openid-configuration"
	response, err := http.Get(well_known_endpoint)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()
	code := response.Status
	fmt.Println(code)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(body, &wellknown_configuration)

}

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetches configuration from well-known endpoint",
	Long: `Sends a GET request to /well-known/.openid-configuration
and fetches the OpenID connect configuration.`,
	Run: fetch,
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
