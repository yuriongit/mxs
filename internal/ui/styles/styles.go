package styles

import (
	lg "github.com/charmbracelet/lipgloss"
)

// Lip Gloss Styles
var (
	Subtle = lg.NewStyle().Faint(true).Italic(true)

	// Purple script name header
	ScriptName = lg.NewStyle().
			Foreground(lg.Color(Purple)).
			Bold(true)

	// Base output box style (border color is dynamic per state)
	BaseOutputBox = lg.NewStyle().
			Border(lg.NormalBorder()).
			Italic(true).
			Faint(true).
			UnsetBorderLeft().
			UnsetBorderRight().
			UnsetBorderTop()
			// UnsetBorderBottom().
			// Padding(0, 1)

	Success = lg.NewStyle().
		Foreground(lg.Color(Green)).
		Bold(true)

	Error = lg.NewStyle().
		Foreground(lg.Color(Red)).
		Bold(true)

	// Purple spinner style
	Spinner = lg.NewStyle().Foreground(lg.Color(Purple))
)

