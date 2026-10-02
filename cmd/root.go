/*
Package cmd holds all of Xs's commands.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>
*/
package cmd

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/yuriongit/xs/internal/ui/execUI"
)

var rootCmd = &cobra.Command{
	Use:   "xs <script-name> [args...]",
	Short: "A CLI tool with AI capabilities for managing and executing scripts",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		scriptName := args[0]
		scriptArgs := args[1:]

		p := tea.NewProgram(execUI.InitialModel(scriptName, scriptArgs))
		if _, err := p.Run(); err != nil {
			log.Fatal("Fatal error running xs", "err", err)
		}
	},
}

func Execute() {
	log.SetPrefix("xs")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
