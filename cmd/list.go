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

func init() {
	var new_config_map map[string]openid.Configuration

	listFunc := func(cmd *cobra.Command, args []string) error {
		b, err := os.ReadFile(config_file)
		if errors.Is(err, os.ErrNotExist) {
			log.Println("Configuration file not found, please run \"loki setup\".")
			return err
		}
		err = json.Unmarshal(b, &new_config_map)
		if err != nil {
			return err
		}
		i := 1
		for _, entry := range new_config_map {
			fmt.Printf("%d: %v, %v\n", i, entry.OpenID.Hostname, entry.Client_id)
			i += 1
		}
		return nil
	}
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List aliases configured",
		Long: `
Lists the hostname and Client ID of aliases configured,`,
		RunE: listFunc,
	}
	rootCmd.AddCommand(listCmd)

}
