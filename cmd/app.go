/*
Package cmd holds all of TDay's commands.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cmd

import "github.com/yuriongit/xs/internal/app"


var application *app.App

// SetApp sets the application struct.
func SetApp(a *app.App) {
	application = a
}

// GetApp gets the application struct.
func GetApp() *app.App {
	return application
}