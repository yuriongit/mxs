/*
Package cnf provides all the functionality
for MXS' config. Manages the application's
configuration and script directory.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cnf

import (
	"errors"
	"strings"
)

// ExtractScriptDescription extracts the description in a user's script.
func (c *Cnf) ExtractScriptDescription(file string) (string, error) {
	_, remaining, found := strings.Cut(file, c.ScriptDescription)
	if !found {
		return "", errors.New("script description not found")
	}

	remaining = strings.TrimLeft(remaining, "\r\n")
	desc, _, _ := strings.Cut(remaining, "\n")
	desc = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(desc), "#"))

	if desc == "" {
		return "", errors.New("script description is empty")
	}

	return desc, nil
}
