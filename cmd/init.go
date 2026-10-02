/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yuriongit/xs/internal/app"
)
var globalApp *app.App

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
  PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		// Skip app init for "init" command
		if cmd.Name() == "init" {
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
	  return globalApp.Cnf.InitCnfDir()
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
