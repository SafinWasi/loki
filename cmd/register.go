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
	var alias string
	var hostname string
	var ssa_file string

	registerFunc := func(cmd *cobra.Command, args []string) error {
		_, exists := new_config_map[alias]
		if exists {
			log.Println("Alias exists. Skipping...")
			os.Exit(1)
		}
		ssa := ""
		if ssa_file != "" {
			b, err := os.ReadFile(ssa_file)
			if err != nil {
				return err
			}
			ssa = string(b)
		}
		b, err := os.ReadFile(config_file)
		if errors.Is(err, os.ErrNotExist) {
			log.Println("Configuration file not found, please run \"loki setup\".")
			return err
		}
		err = json.Unmarshal(b, &new_config_map)
		if err != nil {
			return err
		}
		oidc, err := openid.Register(hostname, ssa)
		if err != nil {
			log.Println("Failed to register client: ", err)
			return err
		}
		new_config_map[alias] = *oidc
		new_bytes, _ := json.MarshalIndent(new_config_map, "", "\t")
		os.WriteFile(config_file, new_bytes, 0644)
		fmt.Println("Registration successful")
		return nil
	}
	var registerCmd = &cobra.Command{
		Use:   "register",
		Short: "Registers a new OpenID Client",
		Long: `
Lists the hostname and Client ID of aliases configured,`,
		RunE: registerFunc,
	}
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringVar(&hostname, "hostname", "", "Hostname of the server")
	registerCmd.Flags().StringVarP(&alias, "alias", "a", "", "Alias of the new client")
	registerCmd.Flags().StringVarP(&ssa_file, "ssa-file", "s", "", "File containing software statement")
	registerCmd.MarkFlagRequired("hostname")
	registerCmd.MarkFlagRequired("alias")
}
