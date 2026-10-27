/*
Package cmd holds all of MXS' commands.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yuriongit/mxs/internal/manager"
)

// addCmd represents the ls command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Creates a new script",
	Long:  "TODO: Implement later",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := GetApp() // TODO: Use GlobalApp instead 
		scriptName := args[0]

		return manager.AddNewScript(app.Cnf, app.Scr, &scriptName)
	},
}

func init() {
	addCmd.SilenceErrors = true
	addCmd.SilenceUsage = true

	rootCmd.AddCommand(addCmd)
}
