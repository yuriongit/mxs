/*
Package ui provides all styles and UI components for MXS.

It provides color constants for Lipgloss styles, colors,
and UI components.
*/
package ui

import (
	"fmt"
	"log"
)

// CreateSuccMsg takes in a success message and returns it with
// a checkmark ([Checkmark]) at the start.
func CreateSuccMsg(msg string) string {
	if msg == "" {
		log.Fatal("success message cannot be empty")
	}

	return fmt.Sprintf("%s %s", Checkmark, msg)
}
