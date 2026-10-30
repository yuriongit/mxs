/*
Package cfg provides all the functionality
for MXS' config. Manages the application's
configuration and script directory.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cfg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const lastExecutedFileName = "last-executed.json"

// LastExecuted maps a script name to its last execution time.
type LastExecuted map[string]time.Time

// LastExecutedPath returns the full path to the last-executed history file.
func (c *Cfg) LastExecutedPath() string {
	return filepath.Join(c.FullCfgPath, lastExecutedFileName)
}

// LoadLastExecuted reads the last-executed history.
//
// A missing history file is treated as empty history.
func (c *Cfg) LoadLastExecuted() (LastExecuted, error) {
	// #nosec G304 -- path built from the config directory and a fixed filename
	raw, err := os.ReadFile(c.LastExecutedPath())
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return LastExecuted{}, nil
		}
		return nil, fmt.Errorf("reading last executed history: %w", err)
	}

	history := LastExecuted{}
	if err := json.Unmarshal(raw, &history); err != nil {
		return nil, fmt.Errorf("decoding last executed history: %w", err)
	}

	return history, nil
}

// RecordLastExecuted records t as the last execution time for scriptName.
func (c *Cfg) RecordLastExecuted(scriptName string, t time.Time) error {
	history, err := c.LoadLastExecuted()
	if err != nil {
		return err
	}

	history[scriptName] = t

	raw, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding last executed history: %w", err)
	}

	// #nosec G304 -- path built from the config directory and a fixed filename
	if err := os.WriteFile(c.LastExecutedPath(), append(raw, '\n'), 0600); err != nil {
		return fmt.Errorf("writing last executed history: %w", err)
	}

	return nil
}
