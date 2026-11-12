/*
Package cmd holds all of MXS' commands.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yuriongit/mxs/internal/scr"
)

// openCmd represents the ls command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Opens a new script",
	Long:  "TODO: Implement later",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		scriptName := args[0]

		return scr.OpenScript(globalApp.Validator, globalApp.Cfg, scriptName)
	},
}

func init() {
	openCmd.SilenceErrors = true
	openCmd.SilenceUsage = true

	rootCmd.AddCommand(openCmd)
}
