/*
Package ui provides all UI components for MXS.

It provides color constants for Lipgloss styles,
colors, and UI components.
*/
package ui

import (
	"errors"
	"fmt"

	lg "github.com/charmbracelet/lipgloss"
)

func CreateSuccessMessage(msg string) (string, error) {
	if msg == "" {
		return "", errors.New("success message cannot be empty")
	}

	coloredCheckMark := Success.Render("✓")

	return fmt.Sprintf("%s %s", coloredCheckMark, msg), nil
}

func CreateInfoHeader(header string, s lg.Style, wall bool) string {
	pipe := ""
	if wall {
		pipe = s.Render("│ ")
	}
	return fmt.Sprintf(
		"%s%s\n%s",
		// TODO: migrate [style] parameter to color once dark/light mode
		s.Bold(true).Render(header),
		s.Faint(true).Render("="),
		pipe,
	)
}
