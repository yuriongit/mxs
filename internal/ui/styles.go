/*
Package ui provides all styles and UI components for MXS.

It provides color constants for Lipgloss styles, colors,
and UI components.
*/
package ui

import (
	lg "charm.land/lipgloss/v2"
)

// Lipgloss Styles
var (
	// Log is the style for scripts logs.
	Log = lg.NewStyle()
	// Primary is a regular style that uses [LightPink].
	Primary = lg.NewStyle().
		Foreground(Magenta)
	// Subtle is just a fainted regular style.
	Subtle = lg.NewStyle().
		Faint(true)
	// Purple script name header.
	ScrName = lg.NewStyle().
		Foreground(Blue)

	// Succ is a regular style that uses [Green].
	Succ = lg.NewStyle().
		Foreground(Green)
	// Err is a regular style that uses [Red].
	Err = lg.NewStyle().
		Foreground(Red)

	// BaseBox is the base style for a box. Includes a left
	// padding of 1.
	BaseBox = lg.NewStyle().
		Border(lg.NormalBorder()).
		UnsetBorderRight().
		UnsetBorderBottom().
		UnsetBorderTop().
		PaddingLeft(1)
	// ExecBox is [BaseBox] with a [Blue] foreground (includes border).
	ExecBox = BaseBox.
		BorderForeground(Blue).
		Foreground(Blue)
	// SuccBox is [BaseBox] with a [Green] foreground (includes border).
	SuccBox = BaseBox.
		BorderForeground(Green).
		Foreground(Green)
	// ErrBox is [BaseBox] with a [Red] foreground (includes border).
	ErrBox = BaseBox.
		BorderForeground(Red).
		Foreground(Red)

	// Primary spinner style that uses [Blue]
	Spinner = lg.NewStyle().Foreground(Blue)
)
