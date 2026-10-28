// Package scr manages script data, creation, validation, and related operations.
package scr

import (
	_ "embed"
)

// Scr holds the script template and script description header.
type Scr struct {
	Template   string
	DescHeader string
}

// templateFile is the embedded script template.
//
//go:embed template.sh
var templateFile string

// New instantiates a new [Scr].
func New() *Scr {
	return &Scr{
		Template:   templateFile,
		DescHeader: "# mxs-description\n# ",
	}
}
