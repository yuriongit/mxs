/*
Package cmd holds all of MXS' commands.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yuriongit/mxs/internal/scr"
)

// laCmd represents the la command
var laCmd = &cobra.Command{
	Use:   "la",
	Short: "Lists all scripts with verbose output",
	Long:  "TODO: Implement later",
	RunE: func(_ *cobra.Command, _ []string) error {
		return scr.ListAllScripts(globalApp.Cfg, true) // true = verbose
	},
}

func init() {
	laCmd.SilenceErrors = true
	laCmd.SilenceUsage = true

	rootCmd.AddCommand(laCmd)
}
