/*
Package ui provides all styles and UI components for MXS.

It provides color constants for Lipgloss styles, colors,
and UI components.
*/
package ui

import (
	"fmt"
	"log"
	"strings"

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

// CreateInfoHeader creates a header with an equals symbol and an optional
// wall (pipe symbol).
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

// ColumnPair represents a two-column row with individual styles for each cell.
//
// LeftPrefix is an optional unstyled-for-width segment rendered before Left
// with its own style (e.g. a metadata tag).
type ColumnPair struct {
	Left            string
	Right           string
	LeftStyle       lg.Style
	LeftPrefix      string
	LeftPrefixStyle lg.Style
}

// FormatAlignedColumns aligns two-column pairs dynamically based on the longest key.
// gap adds extra space between the left and right columns.
func FormatAlignedColumns(rows []ColumnPair, gap int) string {
	maxLen := 0
	for _, row := range rows {
		if width := len(row.LeftPrefix) + len(row.Left); width > maxLen {
			maxLen = width
		}
	}

	var strBuilder strings.Builder
	for idx, row := range rows {
		if idx > 0 {
			strBuilder.WriteString("\n")
		}

		if row.Right != "" {
			// Apply dynamic width to left style to align column 2
			leftStyle := row.LeftStyle.Width(maxLen + gap - len(row.LeftPrefix))
			fmt.Fprintf(
				&strBuilder,
				"%s%s%s",
				row.LeftPrefixStyle.Render(row.LeftPrefix),
				leftStyle.Render(row.Left),
				row.Right,
			)
		} else {
			fmt.Fprintf(
				&strBuilder,
				"%s%s",
				row.LeftPrefixStyle.Render(row.LeftPrefix),
				row.LeftStyle.Render(row.Left),
			)
		}
	}

	return strBuilder.String()
}
