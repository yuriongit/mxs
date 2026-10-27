package scr

import (
	"errors"
	"strings"
)

// ExtractDesc extracts a script's description.
func (s *Scr) ExtractDesc(file string) (string, error) {
	_, remaining, found := strings.Cut(file, s.DescHeader)
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
