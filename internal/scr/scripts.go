// Package scr manages script data, creation, validation, and related operations.
package scr

import (
	_ "embed"
)

type Scr struct {
	Template   string
	DescHeader string
}

// Script template file embedding
//
//go:embed template.sh
var templateFile string

func New() *Scr {
	return &Scr{
		Template:   templateFile,
		DescHeader: "# mxs-description\n# ",
	}
}
