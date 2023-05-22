/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var force bool
var openid_config string

func fetch(cmd *cobra.Command, args []string) {
	Fetch_openid(args[0])
}

func Fetch_openid(hostname string) []byte {
	var old_config map[string]any
	var new_config map[string]any

	b, err := os.ReadFile(openid_config)
	if err != nil {
		log.Println(err)
	}
	if err = json.Unmarshal(b, &old_config); err != nil {
		log.Println(err)
	}
	if old_config[hostname] == nil || force {
		response, err := http.Get(hostname + "/.well-known/openid-configuration")
		if err != nil {
			log.Println(err)
			return nil
		}
		defer response.Body.Close()
		log.Println(response.Status)
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println(err)
			return nil
		}
		if err = json.Unmarshal(body, &new_config); err != nil {
			log.Println(err)
		}
		old_config[hostname] = new_config
		b, _ := json.MarshalIndent(old_config, "", "\t")
		os.WriteFile(openid_config, b, 0644)
	} else {
		log.Println("OpenID configuration exists. Ignoring...")
	}
	return b
}

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetches configuration from well-known endpoint",
	Long: `
Sends a GET request to /well-known/.openid-configuration
and fetches the OpenID connect configuration, then stores it in
the configuration file.`,
	Args: cobra.MatchAll(cobra.ExactArgs(1)),
	Run:  fetch,
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	fetchCmd.PersistentFlags().StringVar(&openid_config, "openid", cwd+string(os.PathSeparator)+".openid.json", "Data file to cache openid configuration entries")
	empty_map := []byte("{}")
	_, err = os.ReadFile(openid_config)
	if errors.Is(err, os.ErrNotExist) {
		log.Println(err)
		log.Println("Creating OpenID file...")
		os.WriteFile(openid_config, empty_map, 0644)
	}
	fetchCmd.Flags().BoolVarP(&force, "force", "f", false, "force retrieval for cached entries")
}
