/*
Package cmd holds all of MXS' commands.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yuriongit/mxs/internal/scr"
)

// forceDelete bypasses the confirmation prompt when removing a script.
var forceDelete bool

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm <script-name>",
	Short: "Removes a script and its metadata",
	Long:  "Removes the script file for <script-name> and its entry in the metadata store.",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return scr.RemoveScript(globalApp.Cfg, args[0], forceDelete, os.Stdin)
	},
}

func init() {
	rmCmd.SilenceErrors = true
	rmCmd.SilenceUsage = true

	rmCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "skip the confirmation prompt")

	rootCmd.AddCommand(rmCmd)
}
