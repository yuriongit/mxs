/*
Package cmd holds all of MXS' commands.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes MXS' config",
	Long:  "TODO: Implement later",
	RunE: func(_ *cobra.Command, _ []string) error {
		ptrApp := GetApp()

		return ptrApp.Cnf.InitCnfDir()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
