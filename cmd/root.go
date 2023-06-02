/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands

var rootCmd = rootGenerator()

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
			_, err := os.ReadFile(config_file)
			if errors.Is(err, os.ErrNotExist) {
				log.Println("Configuration file not found, please run \"loki setup\".")
			}
		},
	}
	rootCmd.Flags().StringVarP(&flow, "flow", "f", "", "Flow to be used for authentication")
	rootCmd.Flags().StringVarP(&alias, "alias", "a", "", "Flow to be used for authentication")
	rootCmd.MarkFlagRequired("alias")
	rootCmd.MarkFlagRequired("flow")
	return rootCmd
}

func init() {

}
