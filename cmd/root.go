/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/safinwasi/loki/openid"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands

var rootCmd = rootGenerator()
var disable_ssl bool

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func rootGenerator() *cobra.Command {
	var flow string
	var alias string
	var rootCmd = &cobra.Command{
		Use:   "loki [command] [flags]",
		Short: "A simple CLI based OAuth RP",
		Long: `
Loki is a simple CLI based application for interacting
with an OAuth 2.0 server. It supports OpenID connect and authorization 
flows. Written in Go.`,
		Run: func(cmd *cobra.Command, args []string) {
			if disable_ssl {
				log.Println("DISABLING SSL (NOT RECOMMENDED)")
			}
			_, err := os.ReadFile(config_file)
			if errors.Is(err, os.ErrNotExist) {
				log.Println("Configuration file not found, please run \"loki setup\".")
				os.Exit(1)
			}
			var config_map map[string]openid.Configuration
			b, err := os.ReadFile(config_file)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			json.Unmarshal(b, &config_map)
			config, exists := config_map[alias]
			if !exists {
				log.Println("Alias missing, please run \"loki setup\".")
				os.Exit(1)
			}
			at, err := openid.Authenticate(flow, config, disable_ssl)
			if err != nil {
				log.Println("Authentication failed:", err)
				os.Exit(1)
			}
			fmt.Printf("Access token obtained: %v\n", at)
		},
	}
	rootCmd.Flags().StringVarP(&flow, "flow", "f", "", "Flow to be used for authentication")
	rootCmd.Flags().StringVarP(&alias, "alias", "a", "", "Flow to be used for authentication")
	rootCmd.PersistentFlags().BoolVar(&disable_ssl, "disable-ssl", false, "Disables SSL")
	rootCmd.MarkFlagRequired("alias")
	rootCmd.MarkFlagRequired("flow")
	return rootCmd
}

func init() {

}
