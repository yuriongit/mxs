/*
Package styles provides all UI components for Xs.

It provides color constants for Lipgloss styles,
colors, and UI components.
*/
package styles

import (
	lg "github.com/charmbracelet/lipgloss"
)

// Color constants.
const (
  // Lime color
	Lime   = "#CCFF5E"
	// Green color
	Green  = "#4CD100"
	// Purple color
	Purple = "#915EFF"
	// Light pink color
	LightPink = "#fb84b5"
	// Pink color
	Pink   = "#FF5ECC"
	// Red color
	Red    = "#FF5E7C"
)

// ============================================================================
// Lip Gloss Styles
// ============================================================================

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
