/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "loki [command] [flags]",
	Short: "A simple CLI based OAuth RP",
	Long: `
Loki is a simple CLI based application for interacting
with an OAuth 2.0 server. It supports OpenID connect and authorization 
flows. Written in Go.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
