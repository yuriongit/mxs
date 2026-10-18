package ui

import (
	lg "github.com/charmbracelet/lipgloss"
)

// Lip Gloss Styles
var (
	Subtle = lg.NewStyle().Faint(true).Italic(true)

	// Purple script name header
	ScriptName = lg.NewStyle().
			Foreground(lg.Color(Pink))

	// Base output box style (border color is dynamic per state)
	BaseOutputBox = lg.NewStyle().
			Border(lg.NormalBorder()).
			UnsetBorderLeft().
			UnsetBorderRight().
			UnsetBorderTop()
		// UnsetBorderBottom().
		// Padding(0, 1)

	Success = lg.NewStyle().
		Foreground(lg.Color(Lime))

	Error = lg.NewStyle().
		Foreground(lg.Color(Red))

	// Purple spinner style
	Spinner = lg.NewStyle().Foreground(lg.Color(Pink))

	DirEntrySymbol = lg.NewStyle().
			Foreground(lg.Color(Blue)).
			Faint(true).
			PaddingLeft(1).
			Render("├─")
	DirEndingEntrySymbol = lg.NewStyle().
				Foreground(lg.Color(Blue)).
				Faint(true).
				PaddingLeft(1).
				Render("└─")
	ScriptEntryStyle = lg.NewStyle().
				PaddingRight(1).
				Foreground(lg.Color(Blue))
)
