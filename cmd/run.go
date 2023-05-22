/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var hostname string
var current_alias string

func authorization_code() error {
	log.Println("Starting authorization code flow")
	return nil
}

func client_credentials() error {
	log.Println("Starting client credentials flow")
	return nil
}

func run(cmd *cobra.Command, args []string) {
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
