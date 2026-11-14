/*
Package cmd holds all of MXS' commands.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cfgCmd represents the ls command
var cfgCmd = &cobra.Command{
	Use:   "cfg",
	Short: "Opens configuration file",
	Long:  "Opens the MXS configuration file in the configured editor",
	RunE: func(_ *cobra.Command, _ []string) error {
		msg, err := globalApp.Cfg.File.Open(globalApp.Cfg, globalApp.Validator)
		if err != nil {
			return err
		}
		fmt.Println(msg)
		return nil
	},
}

func init() {
	cfgCmd.SilenceErrors = true
	cfgCmd.SilenceUsage = true

	rootCmd.AddCommand(cfgCmd)
}
