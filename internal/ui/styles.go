/*
Package ui provides all styles and UI components for MXS.

It provides color constants for Lipgloss styles, colors,
and UI components.
*/
package ui

import (
	lg "charm.land/lipgloss/v2"
)

// Miscellaneous styles.
var (
	// Default style.
	Default = lg.NewStyle()

	// DefaultSubtle is the [Default] style but fainted.
	DefaultSubtle = lg.NewStyle().
			Faint(true)
)

// State styles.
var (
	// Exec is the [Default] style with a [Blue] foreground.
	Exec = Default.Foreground(Blue)

	// Succ is the [Default] style with a [Green] foreground.
	Succ = lg.NewStyle().
		Foreground(Green)
	// Err is the [Default] style with a [Red] foreground.
	Err = lg.NewStyle().
		Foreground(Red)
	// Primary spinner style with a [Blue] foreground.
	Spinner = lg.NewStyle().Foreground(Blue)
)

// Scheme styles.
var (
	// Primary is the [Default] style with a [Magenta] foreground.
	Primary = Default.
		Foreground(Magenta)

	// ScriptReference is the style for referencing scripts; uses a [Blue] foreground.
	ScriptReference = lg.NewStyle().
			Foreground(Blue)
)

// Background styles.
var (
	// highlighted should have a background when being used.
	highlighted = Default.Padding(0, 1, 0, 1)

	// ConfigReference is "config" highlighted with a [Magenta] background (uses [highlighted] style).
	// It's purpose is to signify that "config" is not a file, but a speciaL reference specifically
	// mentioning the user's actual configuration file (no other file will receive any type of
	// highlighted reference obviously).
	//
	// Usage: e.g., "opened 'config' with vim"
	ConfigReference = highlighted.Foreground(lg.Color("15")).Background(Magenta).Render("config")
)

// Box styles.
var (
	// BaseBox is the base style for a box.
	BaseBox = lg.NewStyle().
		Border(lg.NormalBorder()).
		UnsetBorderRight().
		UnsetBorderBottom().
		UnsetBorderTop().
		PaddingLeft(1)

	// ExecBox is the [BaseBox] style with a [Blue] foreground.
	ExecBox = BaseBox.
		BorderForeground(Blue).
		Foreground(Blue)

	// SuccBox is the [BaseBox] style with a [Green] foreground.
	SuccBox = BaseBox.
		BorderForeground(Green).
		Foreground(Green)

	// ErrBox is the [BaseBox] style with a [Red] foreground.
	ErrBox = BaseBox.
		BorderForeground(Red).
		Foreground(Red)
)
