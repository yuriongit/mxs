package ui

import (
	lg "github.com/charmbracelet/lipgloss"
)

// Lip Gloss Styles
var (
	// Scripts outputted logs style
	SubMsgStyle = lg.NewStyle().Faint(true)
	// Script description and outro message style
	ScriptDescAndOutroMsgStyle = lg.NewStyle().Foreground(lg.Color(Blue))
	// Regular text style (dark mode supported)
	RegularTextStyle = lg.NewStyle().Foreground(lg.Color(RegularTextColor))

	Subtle = lg.NewStyle().Faint(true)

	// Purple script name header
	ScriptName = lg.NewStyle().
			Foreground(lg.Color(Pink))

	// Base output box style (border color is dynamic per state)
	BaseOutputBox = lg.NewStyle().
			Border(lg.NormalBorder()).
			UnsetBorderRight().
			UnsetBorderBottom().
			UnsetBorderTop().
			PaddingLeft(1)

	RunningBoxStyle = BaseOutputBox.BorderForeground(lg.Color(Pink))

	ScriptsBoxStyle = BaseOutputBox.BorderForeground(lg.Color(Blue))

	SuccessBoxStyle = BaseOutputBox.BorderForeground(lg.Color(Lime)).Foreground(lg.Color(Lime))
	Success         = lg.NewStyle().
			Foreground(lg.Color(Lime))

	Error = lg.NewStyle().
		Foreground(lg.Color(Red))
	ErrBoxStyle = BaseOutputBox.BorderForeground(lg.Color(Red))

	// Purple spinner style
	Spinner = lg.NewStyle().Foreground(lg.Color(Pink))

	ScriptEntryStyle = lg.NewStyle().
				PaddingRight(1).
				Foreground(lg.Color(Blue)).
				Faint(true)

	ScriptsDirStyle = lg.NewStyle().
			Foreground(lg.Color(Blue)).
			PaddingLeft(0).
			Bold(true)
)
