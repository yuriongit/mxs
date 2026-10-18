package ui

import (
	lg "github.com/charmbracelet/lipgloss"
)

var (
	Pipe = "│"
)

var (
	SubMsgStyle                = lg.NewStyle().Faint(true).PaddingLeft(2)
	ScriptDescAndOutroMsgStyle = lg.NewStyle().Foreground(lg.Color(Blue))
	RegularTextStyle           = lg.NewStyle().Foreground(lg.Color(RegularTextColor))
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
	RunningBoxStyle = BaseOutputBox.BorderForeground(lg.Color(Pink))

	SuccessBoxStyle = BaseOutputBox.BorderForeground(lg.Color(Lime)).Foreground(lg.Color(Lime))
	Success         = lg.NewStyle().
			Foreground(lg.Color(Lime))

	Error = lg.NewStyle().
		Foreground(lg.Color(Red))
	ErrBoxStyle = BaseOutputBox.BorderForeground(lg.Color(Red))

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

	ScriptsDirStyle = lg.NewStyle().
			Foreground(lg.Color(Pink)).
			PaddingLeft(0).
			Bold(true)
)
