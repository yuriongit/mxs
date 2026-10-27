/*
Package ui provides all styles and UI components for MXS.

It provides color constants for Lipgloss styles, colors,
and UI components.
*/
package ui

import (
	lg "github.com/charmbracelet/lipgloss"
)

// Lipgloss Styles
var (
	// Log is the style for scripts logs.
	Log = lg.NewStyle().
		Faint(true)
	// Primary is a regular style that uses [Blue].
	Primary = lg.NewStyle().
		Foreground(lg.Color(Blue))
	// Subtle is just a fainted regular style.
	Subtle = lg.NewStyle().
		Faint(true)
	// Purple script name header.
	ScrName = lg.NewStyle().
		Foreground(lg.Color(Pink))

	// Succ is a regular style that uses [Lime].
	Succ = lg.NewStyle().
		Foreground(lg.Color(Lime))
	// Err is a regular style that uses [Red].
	Err = lg.NewStyle().
		Foreground(lg.Color(Red))

	// BaseBox is the base style for a box. Includes a left
	// padding of 1.
	BaseBox = lg.NewStyle().
		Border(lg.NormalBorder()).
		UnsetBorderRight().
		UnsetBorderBottom().
		UnsetBorderTop().
		PaddingLeft(1)
	// ExecBox is [BaseBox] with a [Pink] border foreground.
	ExecBox = BaseBox.BorderForeground(lg.Color(Pink))
	// SuccBox is [BaseBox] with a [Lime] border foreground.
	SuccBox = BaseBox.
		BorderForeground(lg.Color(Lime)).
		Foreground(lg.Color(Lime))
	// ErrBox is [BaseBox] with a [Red] border foreground.
	ErrBox = BaseBox.BorderForeground(lg.Color(Red))

	// Primary spinner style that uses [Pink]
	Spinner = lg.NewStyle().Foreground(lg.Color(Pink))

	// DirEntrySymbol styles a branch connector.
	DirEntrySymbol = lg.NewStyle().
			Foreground(lg.Color(Blue)).
			Faint(true).
			PaddingLeft(1).
			Render("├─")
	// DirEndingEntrySymbol styles the last branch connector.
	DirEndingEntrySymbol = lg.NewStyle().
				Foreground(lg.Color(Blue)).
				Faint(true).
				PaddingLeft(1).
				Render("└─")
	// ScriptEntryStyle styles a script entry.
	ScriptEntryStyle = lg.NewStyle().
				PaddingRight(1).
				Foreground(lg.Color(Blue))
	// ScriptsDirStyle styles the scripts directory label.
	ScriptsDirStyle = lg.NewStyle().
			Foreground(lg.Color(Pink)).
			PaddingLeft(0).
			Bold(true)
)
