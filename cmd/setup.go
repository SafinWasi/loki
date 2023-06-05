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

var config_file string

func init() {
	var new_config openid.Configuration
	var new_config_map map[string]openid.Configuration
	var alias string
	var new_bytes []byte

	setupFunc := func(cmd *cobra.Command, args []string) error {
		old_bytes, err := os.ReadFile(config_file)
		if errors.Is(err, os.ErrNotExist) {
			new_config_map = make(map[string]openid.Configuration)
		} else {
			json.Unmarshal(old_bytes, &new_config_map)
			_, exists := new_config_map[alias]
			if exists {
				log.Println("Alias exists. Skipping...")
				return nil
			}
		}
		oidc, err := openid.Fetch_openid(new_config.OpenID.Hostname, disable_ssl)
		if err != nil {
			return err
		}
		new_config.OpenID = *oidc
		new_config_map[alias] = new_config
		new_bytes, _ = json.MarshalIndent(new_config_map, "", "\t")
		os.WriteFile(config_file, new_bytes, 0644)
		fmt.Println("Setup successful")
		return nil
	}
	var setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Setup details for OIDC",
		Long: `
This command must be run before authentication attempt can be made,`,
		RunE: setupFunc,
	}
	rootCmd.AddCommand(setupCmd)
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	filename := pwd + string(os.PathSeparator) + ".config.json"
	setupCmd.PersistentFlags().StringVar(&config_file, "config", filename, "Location of configuration file")
	setupCmd.Flags().StringVar(&new_config.OpenID.Hostname, "hostname", "http://127.0.0.1", "Hostname of OpenID Provider")
	setupCmd.Flags().StringVar(&new_config.Client_id, "client-id", "", "Client ID")
	setupCmd.Flags().StringVar(&new_config.Client_secret, "client-secret", "", "Client Secret")
	setupCmd.Flags().StringVar(&alias, "alias", "", "Alias for configuration")
	setupCmd.MarkFlagRequired("client-id")
	setupCmd.MarkFlagRequired("client-secret")
	setupCmd.MarkFlagRequired("alias")
}
