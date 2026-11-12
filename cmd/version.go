/*
Package cmd holds all of MXS' commands.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the ls command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version of MXS",
	Long:  "TODO: Implement later",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println(globalApp.Version)
	},
}

func init() {
	versionCmd.SilenceErrors = true
	versionCmd.SilenceUsage = true

	rootCmd.AddCommand(versionCmd)
}
