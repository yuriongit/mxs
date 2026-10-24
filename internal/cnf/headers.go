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

// ExtractScriptDescHeader extracts the description in a user's script.
func (c *Cnf) ExtractScriptDescHeader(file string) (string, error) {
	_, remaining, found := strings.Cut(file, c.ScriptDescHeader)
	if !found {
		return "", errors.New("script description header not found")
	}

	remaining = strings.TrimLeft(remaining, "\r\n")
	line, _, _ := strings.Cut(remaining, "\n")
	line = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "#"))

	if line == "" {
		return "", errors.New("script description is empty")
	}

	return line, nil
}

// ValidateScript validates that required script headers exist.
func (c *Cnf) ValidateScript(file string) (scriptDesc string, err error) {
	return c.ExtractScriptDescHeader(file)
}
