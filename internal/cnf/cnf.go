/*
Package cnf provides all the functionality
for MXS' config. Manages the application's
configuration and script directory.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cnf

import (
	"errors"
	"os"
	"path/filepath"
)

type Cnf struct {
	Dir             string
	ScriptsPath     string
	FullCnfPath     string
	FullScriptsPath string
	ScriptTemplate  string
	File            *File
	FileName        string
}

var demoScriptFileName = "demo.sh"

func NewCnf(homePath string) *Cnf {
	cnfDir := ".mxs"
	scriptsDir := filepath.Join(cnfDir, "scripts")
	fullCnfPath := filepath.Join(homePath, cnfDir)
	fullScriptsPath := filepath.Join(homePath, scriptsDir)
	scriptTemplate := `#!mxs-description
# description

echo "%s"

#!mxs-end
# done`
	cnfFileName := "mxs.json"

	return &Cnf{
		cnfDir,
		scriptsDir,
		fullCnfPath,
		fullScriptsPath,
		scriptTemplate,
		&File{},
		cnfFileName,
	}
}

func (c *Cnf) hasCnfDirectory(subDir string) error {
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
