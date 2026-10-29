/*
Package cfg provides all the functionality
for MXS' config. Manages the application's
configuration and script directory.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cfg

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yuriongit/mxs/internal/ui"
)

// InitCfgDir creates the configuration and scripts directories if they
// don't already exist.
func (c *Cfg) InitCfgDir() error {
	directories := []struct {
		path string
		name string
	}{
		{path: c.FullCfgPath, name: fmt.Sprintf("~/%s", c.Dir)},
		{path: c.FullScriptsPath, name: fmt.Sprintf("~/%s", c.ScriptsPath)},
	}

	for _, directory := range directories {
		info, err := os.Stat(directory.path)

		switch {
		case err == nil:
			if !info.IsDir() {
				return fmt.Errorf("%s exists but is not a directory", directory.name)
			}
			fmt.Println(ui.CreateSuccMsg(fmt.Sprintf("found existing %s directory", directory.name)))
		case os.IsNotExist(err):
			if err := os.MkdirAll(directory.path, 0700); err != nil {
				return fmt.Errorf(
					"failed to create %s directory: %w",
					directory.name,
					err,
				)
			}
			fmt.Println(ui.CreateSuccMsg(fmt.Sprintf("created %s directory", directory.name)))
		default:
			return fmt.Errorf("failed to check %s directory: %w", directory.name, err)
		}
	}

	if err := c.CreateCfgFile(); err != nil {
		return err
	}

	if err := c.SetupDemoScript(); err != nil {
		return err
	}

	return nil
}

//go:embed demo.sh
var demoScriptContent string

// SetupDemoScript creates the demo script.
func (c *Cfg) SetupDemoScript() error {
	fullDemoScriptPath := filepath.Join(c.FullScriptsPath, demoScriptFileName)

	// os.O_EXCL creates the file atomically or fails if it already exists
	// #nosec G304 G302 -- fullDemoScriptPath built from config path with hardcoded filename && enables script execution
	file, err := os.OpenFile(fullDemoScriptPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0700)
	switch {
	case os.IsExist(err):
		fmt.Println(ui.CreateSuccMsg("found existing demo script"))
		return nil
	case err != nil:
		return fmt.Errorf("failed to check or create demo script: %w", err)
	}

	_, writeErr := file.WriteString(demoScriptContent)
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close demo script file: %w", err)
	}
	if writeErr != nil {
		return fmt.Errorf("failed to write demo script: %w", writeErr)
	}
	fmt.Println(ui.CreateSuccMsg("created demo script"))

	return nil
}

//go:embed mxs.json
var cfgFileTemplate string

// CreateCfgFile creates a templatized MXS configurationf file
func (c *Cfg) CreateCfgFile() error {
	fullCfgFilePath := filepath.Join(c.FullCfgPath, c.FileName)

	// os.O_EXCL creates the file atomically or fails if it already exists
	// #nosec G304 G302 -- fullDemoScriptPath built from config path with hardcoded filename && enables script execution
	file, err := os.OpenFile(fullCfgFilePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	switch {
	case os.IsExist(err):
		fmt.Println(ui.CreateSuccMsg("found existing configuration file"))
		return nil
	case err != nil:
		return fmt.Errorf("failed to check or create configuration file: %w", err)
	}

	_, writeErr := file.WriteString(cfgFileTemplate)
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close configuration file: %w", err)
	}
	if writeErr != nil {
		return fmt.Errorf("failed to write configuration file: %w", writeErr)
	}
	fmt.Println(ui.CreateSuccMsg("created configuration file"))

	return nil
}
