/*
Package cnf provides all the functionality
for Xs's config. Manages the application's
configuration and script directory.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cnf

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

func (c *Cnf) InitCnfDir() error {
	directories := []struct {
		path string
		name string
	}{
		{path: c.FullCnfPath, name: fmt.Sprintf("~/%s", c.Dir)},
		{path: c.FullScriptsPath, name: fmt.Sprintf("~/%s", c.ScriptsPath)},
	}

	for _, directory := range directories {
		info, err := os.Stat(directory.path)

		switch {
		case err == nil:
			if !info.IsDir() {
				return fmt.Errorf("%s exists but is not a directory", directory.name)
			}

			fmt.Printf("✓ Found existing %s directory\n", directory.name)

		case os.IsNotExist(err):
			if err := os.MkdirAll(directory.path, 0700); err != nil {
				return fmt.Errorf(
					"Failed to create %s directory: %w",
					directory.name,
					err,
				)
			}

			fmt.Printf("✓ Created %s directory\n", directory.name)

		default:
			return fmt.Errorf("Failed to check %s directory: %w", directory.name, err)
		}
	}

	if err := c.SetupDemoScript(); err != nil {
		return err
	}

	return nil
}

//go:embed demo.sh
var demoScriptContent string

// SetupDemoScript creates the demo script.
func (c *Cnf) SetupDemoScript() error {
	fullDemoScriptPath := filepath.Join(c.FullScriptsPath, demoScriptFileName)

	// os.O_EXCL creates the file atomically or fails if it already exists
	// #nosec G304 -- fullDemoScriptPath built from config path with hardcoded filename
	file, err := os.OpenFile(fullDemoScriptPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	switch {
	case os.IsExist(err):
		fmt.Printf("✓ Found existing demo script in ~/%s\n", c.ScriptsPath)
		return nil
	case err != nil:
		return fmt.Errorf("Failed to check or create demo script: %w", err)
	}

	_, writeErr := file.WriteString(demoScriptContent)
	if err := file.Close(); err != nil {
		return fmt.Errorf("Failed to close demo script file: %w", err)
	}
	if writeErr != nil {
		return fmt.Errorf("Failed to write demo script: %w", writeErr)
	}
	fmt.Printf("✓ Created demo script (%s) in %s\n", demoScriptFileName, c.ScriptsPath)

	return nil
}
