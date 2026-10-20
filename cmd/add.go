/*
Package cmd holds all of MXS' commands.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the ls command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Creates a new script",
	Long:  "TODO: Implement later",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			fmt.Println("Script name:", args[0])
			return nil
			// return manager.AddNewScript(args[0])
		}
		return errors.New("Script name required")
	},
}

func init() {
	addCmd.SilenceErrors = true
	addCmd.SilenceUsage = true

	rootCmd.AddCommand(addCmd)
}
