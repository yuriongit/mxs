/*
Package ui provides all UI components for MXS.

It provides color constants for Lipgloss styles,
colors, and UI components.
*/
package ui

import (
	"errors"
	"fmt"
)

func CreateSuccessMessage(msg string) (string, error) {
	if msg == "" {
		return "", errors.New("success message cannot be empty")
	}

	coloredCheckMark := Success.Render("✓")

	return fmt.Sprintf("%s %s", coloredCheckMark, msg), nil
}
