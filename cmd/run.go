/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/safinwasi/loki/client"
	"github.com/spf13/cobra"
)

var hostname string
var current_alias string
var current_client_id string
var current_client_secret string

func authorization_code() error {
	log.Println("Starting authorization code flow")
	return nil
}

func client_credentials() error {
	log.Println("Starting client credentials flow")
	var oidc_details map[string]any
	b, err := os.ReadFile(openid_config)
	if err != nil {
		return err
	}
	var configuration map[string]any
	json.Unmarshal(b, &configuration)
	if configuration[hostname] == nil {
		log.Println("OpenID configuration missing. Fetching...")
		Fetch_openid(hostname)
		b, err = os.ReadFile(openid_config)
		if err != nil {
			return err
		}
		json.Unmarshal(b, &configuration)
	}
	oidc_details = configuration[hostname].(map[string]any)
	token_endpoint := oidc_details["token_endpoint"]
	concatenated_creds := current_client_id + ":" + current_client_secret
	encoded_creds := base64.RawURLEncoding.EncodeToString([]byte(concatenated_creds))

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "openid")
	r, err := http.NewRequest(http.MethodPost, token_endpoint.(string), strings.NewReader(data.Encode()))
	if err != nil {
		log.Println(err)
		return err
	}
	r.Header.Add("Authorization", "Basic "+encoded_creds)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	log.Println(resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	var formatted map[string]any
	json.Unmarshal(body, &formatted)
	log.Println("Access token obtained")
	return nil
}

func run(cmd *cobra.Command, args []string) {
	var creds client.ClientDict
	b, err := os.ReadFile(client_config)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(b, &creds)
	var current_config client.Client = creds[current_alias]
	current_client_id = current_config.Client_id
	current_client_secret = current_config.Client_secret
	if args[0] == "code" {
		authorization_code()
	} else {
		client_credentials()
	}
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a flow.",
	Long: `
Starts an OAuth 2.0 authentication flow. This may be authorization code or client credentials.
See the flags for more information.`,
	ValidArgs: []string{"code", "credentials"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:       run,
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVar(&hostname, "hostname", "127.0.0.1", "Hostname of authorization server")
	runCmd.Flags().StringVarP(&current_alias, "alias", "a", "", "Alias of client to use in this flow")
	runCmd.MarkFlagRequired("alias")
}
