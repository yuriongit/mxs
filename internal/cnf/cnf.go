/*
Package cnf provides all the functionality
for Xs's config.

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

func NewCnf(homePath string) *Cnf {
	cnfDir := ".xs"
	scriptsDir := "scripts"
	fullCnfPath := filepath.Join(homePath, cnfDir)
	fullScriptsPath := filepath.Join(homePath, cnfDir, scriptsDir)

	return &Cnf{
		cnfDir,
		scriptsDir,
		fullCnfPath,
		fullScriptsPath,
	}
}
