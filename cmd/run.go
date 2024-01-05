package cmd

import (
	"github.com/safinwasi/loki/secrets"
	"github.com/safinwasi/loki/web"
	"github.com/spf13/cobra"
)

var port int

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the Loki server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		secrets.Initialize(Debug)
		web.Start(port)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntVarP(&port, "port", "p", 3000, "Port number to bind to")
}
