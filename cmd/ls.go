/*
Package cmd holds all of MXS' commands.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yuriongit/mxs/internal/scr"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists all scripts",
	Long:  "TODO: Implement later",
	RunE: func(_ *cobra.Command, _ []string) error {
		app := GetApp() // TODO: Use GlobalApp variable

		return scr.ListAllScripts(app.Cnf)
	},
}

func init() {
	lsCmd.SilenceErrors = true
	lsCmd.SilenceUsage = true

	rootCmd.AddCommand(lsCmd)
}
