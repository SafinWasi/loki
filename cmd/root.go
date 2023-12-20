package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var Debug bool

var rootCmd = &cobra.Command{
	Use:   "loki",
	Short: "",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Debug mode")
}
