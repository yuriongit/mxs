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

// metadataFileName is the fixed name of the metadata store inside the config dir.
const metadataFileName = "metadata.json"

// Metadata is the generic, script-scoped metadata store persisted as JSON.
//
// Each entry is keyed by script name; the value holds that script's metadata
// (currently its last execution time). It is intentionally a loose map so new
// metadata kinds can be added without a schema change.
type Metadata map[string]any

// MetadataPath returns the full path to the metadata store file.
func (c *Cfg) MetadataPath() string {
	return filepath.Join(c.FullCfgPath, metadataFileName)
}

// LoadMetadata reads and decodes the metadata store.
//
// A missing store is treated as empty.
func (c *Cfg) LoadMetadata() (Metadata, error) {
	// #nosec G304 -- path built from the config directory and a fixed filename
	raw, err := os.ReadFile(c.MetadataPath())
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Metadata{}, nil
		}
		return nil, fmt.Errorf("reading metadata store: %w", err)
	}

	metadata := Metadata{}
	if err := json.Unmarshal(raw, &metadata); err != nil {
		return nil, fmt.Errorf("decoding metadata store: %w", err)
	}

	return metadata, nil
}

// SaveMetadata writes metadata to the store, creating or overwriting it.
func (c *Cfg) SaveMetadata(metadata Metadata) error {
	raw, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding metadata store: %w", err)
	}

	// #nosec G304 -- path built from the config directory and a fixed filename
	if err := os.WriteFile(c.MetadataPath(), append(raw, '\n'), 0600); err != nil {
		return fmt.Errorf("writing metadata store: %w", err)
	}

	return nil
}

// DeleteMetadataEntry removes the metadata entry for scriptName.
//
// Removing a missing entry is not an error; the store is still written so a
// stale entry is normalized away.
func (c *Cfg) DeleteMetadataEntry(scriptName string) error {
	metadata, err := c.LoadMetadata()
	if err != nil {
		return err
	}

	delete(metadata, scriptName)

	return c.SaveMetadata(metadata)
}

// LastExecuted maps a script name to its last execution time.
//
// It is a typed view over the generic [Metadata] store, kept for the `la`
// command's last-execution display.
type LastExecuted map[string]time.Time

// LoadLastExecuted reads the last-execution metadata.
func (c *Cfg) LoadLastExecuted() (LastExecuted, error) {
	raw, err := c.LoadMetadata()
	if err != nil {
		return nil, err
	}

	lastExecuted := LastExecuted{}
	for name, value := range raw {
		parsed, err := timeValue(value)
		if err != nil {
			return nil, err
		}
		if !parsed.IsZero() {
			lastExecuted[name] = parsed
		}
	}

	return lastExecuted, nil
}

// RecordLastExecuted records t as the last execution time for scriptName.
func (c *Cfg) RecordLastExecuted(scriptName string, t time.Time) error {
	metadata, err := c.LoadMetadata()
	if err != nil {
		return err
	}

	metadata[scriptName] = t

	return c.SaveMetadata(metadata)
}

// timeValue converts a decoded metadata value into a time.Time.
//
// JSON stores times as RFC3339 strings, so this accepts both that form and an
// already-parsed time.
func timeValue(value any) (time.Time, error) {
	switch v := value.(type) {
	case time.Time:
		return v, nil
	case string:
		parsed, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return time.Time{}, fmt.Errorf("decoding stored time %q: %w", v, err)
		}
		return parsed, nil
	default:
		return time.Time{}, fmt.Errorf("metadata value %v is not a time", value)
	}
}
