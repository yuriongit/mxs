/*
Package cmd holds all of MXS' commands.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/yuriongit/mxs/internal/app"
	"github.com/yuriongit/mxs/internal/ui/screens/execui"
)

var globalApp *app.App

var rootCmd = &cobra.Command{
	Use:   "mxs <script-name> [args...]",
	Short: "A CLI tool with AI capabilities for managing and executing scripts",
	Args:  cobra.MinimumNArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		// Skip app init for "init" command
		if cmd.Name() == "help" {
			return nil
		}

		// Initialize app for all other commands
		var err error
		globalApp, err = app.NewApp()
		if err != nil {
			return err
		}

		SetApp(globalApp)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		scriptName := args[0]
		scriptArgs := args[1:]

		p := tea.NewProgram(execui.InitialModel(scriptName, scriptArgs))
		if _, err := p.Run(); err != nil {
			log.Fatal("Fatal error running MXS", "err", err)
		}

		// Initialize app for all other commands
		var err error
		globalApp, err = app.NewApp()
		if err != nil {
			return err
		}

		SetApp(globalApp)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ %s\n", err.Error())
		return
	}

	// Global cleanup
	// if globalApp != nil {
	// 	globalApp.Cancel()
	// 	globalApp.Database.Pool.Close()
	// }
}
