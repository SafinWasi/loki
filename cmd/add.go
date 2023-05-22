/*
Copyright © 2023 Safin Wasi <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/safinwasi/loki/client"
	"github.com/spf13/cobra"
)

var client_config string
var client_id string
var client_secret string
var alias string

func add(cmd *cobra.Command, args []string) {
	clients := make(client.ClientDict)
	clients.AddClient(client_id, client_secret, alias)
	b, err := json.MarshalIndent(clients, "", "\t")
	if err != nil {
		log.Println(err)
	}
	os.WriteFile(client_config, b, 0644)
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds an OAuth2 client",
	Long: `
Saves an OAuth2 client locally for authentication, with an alias for shorthand
If an existing client with the same alias exists, it will be overridden`,
	Run: add,
}

func init() {
	rootCmd.AddCommand(addCmd)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	addCmd.PersistentFlags().StringVar(&client_config, "clientfile", cwd+string(os.PathSeparator)+".client.json", "Data file to store client IDs and secrets")
	empty_map := []byte("{}")
	_, err = os.ReadFile(client_config)
	if errors.Is(err, os.ErrNotExist) {
		log.Println(err)
		log.Println("Creating client file...")
		os.WriteFile(client_config, empty_map, 0644)
	}
	addCmd.Flags().StringVar(&client_id, "id", "", "Client ID")
	addCmd.Flags().StringVar(&client_secret, "secret", "", "Client Secret")
	addCmd.Flags().StringVar(&alias, "alias", "", "Alias for client")
	addCmd.MarkFlagRequired("id")
	addCmd.MarkFlagRequired("secret")
	addCmd.MarkFlagRequired("alias")
}
