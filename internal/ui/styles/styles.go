package styles

import (
	lg "github.com/charmbracelet/lipgloss"
)

// Lip Gloss Styles
var (
	Subtle = lg.NewStyle().Faint(true).Italic(true)

	// Light pink script name header
	ScriptName = lg.NewStyle().
			Foreground(lg.Color(LightPink)).
			Bold(true)

	// Base output box style (border color is dynamic per state)
	BaseOutputBox = lg.NewStyle().
				Border(lg.RoundedBorder()).
				Padding(0, 1).
				MarginTop(1).
				MarginBottom(1)

	Success = lg.NewStyle().
			Foreground(lg.Color(Green)).
			Bold(true)

	Error = lg.NewStyle().
			Foreground(lg.Color(Red)).
			Bold(true)

	// Light pink spinner style
	Spinner = lg.NewStyle().Foreground(lg.Color(LightPink))
)

