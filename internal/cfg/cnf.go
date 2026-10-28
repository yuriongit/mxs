/*
Package cnf provides all the functionality
for MXS' config. Manages the application's
configuration and script directory.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cnf

import (
	_ "embed" // used to embed files into the binary.
	"errors"
	"os"
	"path/filepath"
)

// Cnf manages application configuration and configuration directories.
type Cnf struct {
	Dir             string
	ScriptsPath     string
	FullCnfPath     string
	FullScriptsPath string
	File            *File
	FileName        string
}

var demoScriptFileName = "demo.sh"

// NewCnf instantiates a new [Cnf].
func NewCnf(homePath string) *Cnf {
	cnfDir := ".mxs"
	scriptsDir := filepath.Join(cnfDir, "scripts")
	fullCnfPath := filepath.Join(homePath, cnfDir)
	fullScriptsPath := filepath.Join(homePath, scriptsDir)
	cnfFileName := "mxs.json"

	return &Cnf{
		cnfDir,
		scriptsDir,
		fullCnfPath,
		fullScriptsPath,
		&File{},
		cnfFileName,
	}
}

// HasCnfDir checks if the configuration directory or any sub-directory exists.
//
// Currently, availabe sub-directories are only limited to the scripts folder.
func (c *Cnf) HasCnfDir(subDir string) error {
	switch subDir {
	case "":
		if _, err := os.ReadDir(c.FullCnfPath); err != nil {
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
