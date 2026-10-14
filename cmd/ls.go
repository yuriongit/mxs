/*
Package cmd holds all of MXS' commands.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists all scripts",
	Long:  "TODO: Implement later",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
