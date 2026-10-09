/*
Package cnf provides all the functionality
for MXS' config. Manages the application's
configuration and script directory.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cnf

import (
	"path/filepath"
)

type Cnf struct {
	Dir             string
	ScriptsPath     string
	FullCnfPath     string
	FullScriptsPath string
}

var demoScriptFileName = "demo.sh"

func NewCnf(homePath string) *Cnf {
	cnfDir := ".mxs"
	scriptsDir := filepath.Join(cnfDir, "scripts")
	fullCnfPath := filepath.Join(homePath, cnfDir)
	fullScriptsPath := filepath.Join(homePath, scriptsDir)

	return &Cnf{
		cnfDir,
		scriptsDir,
		fullCnfPath,
		fullScriptsPath,
	}
}
