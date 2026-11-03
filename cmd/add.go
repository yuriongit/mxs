/*
Package cmd holds all of MXS' commands.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yuriongit/mxs/internal/scr"
)

// addCmd represents the ls command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Creates a new script",
	Long:  "TODO: Implement later",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := GetApp() // TODO: Use GlobalApp instead
		scriptName := args[0]

		desc, err := cmd.Flags().GetString("desc")
		if err != nil {
  		return err
		}

		return scr.AddNewScript(app.Cfg, app.Scr, desc, scriptName)
	},
}

func init() {
	addCmd.SilenceErrors = true
	addCmd.SilenceUsage = true

	addCmd.Flags().StringP("desc", "d", "description", "script description")

	rootCmd.AddCommand(addCmd)
}
