/*
Package cfg provides all the functionality
for MXS' config. Manages the application's
configuration and script directory.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cfg

import (
	_ "embed" // used to embed files into the binary.
	"errors"
	"os"
	"path/filepath"
)

// Cfg manages application configuration and configuration directories.
type Cfg struct {
	Dir             string
	ScriptsPath     string
	FullCfgPath     string
	FullScriptsPath string
	File            *File
	FileName        string
}

var demoScriptFileName = "demo.sh"

// NewCfg instantiates a new [Cfg].
func NewCfg(homePath string) *Cfg {
	cfgDir := ".mxs"
	scriptsDir := filepath.Join(cfgDir, "scripts")
	fullCfgPath := filepath.Join(homePath, cfgDir)
	fullScriptsPath := filepath.Join(homePath, scriptsDir)
	cfgFileName := "mxs.json"

	return &Cfg{
		cfgDir,
		scriptsDir,
		fullCfgPath,
		fullScriptsPath,
		&File{},
		cfgFileName,
	}
}

// HasCfgDir checks if the configuration directory or any sub-directory exists.
//
// Currently, availabe sub-directories are only limited to the scripts folder.
func (c *Cfg) HasCfgDir(subDir string) error {
	switch subDir {
	case "":
		if _, err := os.ReadDir(c.FullCfgPath); err != nil {
			return err
		}
	case "scripts":
		if _, err := os.ReadDir(c.FullScriptsPath); err != nil {
			return err
		}
	default:
		return errors.New("invalid sub-directory option")
	}

	return nil
}
