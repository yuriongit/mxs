/*
Package ui provides all styles and UI components for MXS.

It provides color constants for Lipgloss styles, colors,
and UI components.
*/
package ui

import (
	"fmt"
	"log"

	lg "github.com/charmbracelet/lipgloss"
)

// CreateSuccMsg takes in a success message and returns it with
// a checkmark ([Checkmark]) at the start.
func CreateSuccMsg(msg string) string {
	if msg == "" {
		log.Fatal("success message cannot be empty")
	}

	return fmt.Sprintf("%s %s", Checkmark, msg)
}

func CreateInfoHeader(header string, s lg.Style, wall bool) string {
	pipe := ""
	if wall {
		pipe = s.Render("│ ")
	}
	return fmt.Sprintf(
		"%s%s\n%s",
		// TODO: migrate [s (style)] parameter to [lg.Color] once dark/light mode
		// is being integrated
		s.Bold(true).Render(header),
		s.Faint(true).Render("="),
		pipe,
	)
}
