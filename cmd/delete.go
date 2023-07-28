package cmd

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/safinwasi/loki/openid"
	"github.com/spf13/cobra"
)

func init() {
	var new_config_map map[string]openid.Configuration
	var alias string

	deleteFunc := func(cmd *cobra.Command, args []string) error {
		b, err := os.ReadFile(config_file)
		if errors.Is(err, os.ErrNotExist) {
			log.Println("Configuration file not found, please run \"loki setup\".")
			return err
		}
		err = json.Unmarshal(b, &new_config_map)
		if err != nil {
			return err
		}
		delete(new_config_map, alias)
		b, err = json.MarshalIndent(new_config_map, "", "\t")
		if err != nil {
			return err
		}
		err = os.WriteFile(config_file, b, 0644)
		if err == nil {
			log.Printf("Delete %v successful", alias)
		}
		return err
	}
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Deletes a configuration by alias",
		Long: `
Deletes the configuration by the provided alias`,
		RunE: deleteFunc,
	}
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&alias, "alias", "a", "", "alias of configuration to delete")

}
