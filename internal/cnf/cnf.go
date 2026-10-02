/*
Package cnf provides all the functionality
for Xs's config.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cnf

import (
	"fmt"
	"os"
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

func (c *Cnf) ChToBaseDir() error {
	if err := os.Chdir(c.FullCnfPath); err != nil {
		return fmt.Errorf("%w\nfull cnf path: %s", err, c.FullCnfPath)
	}
	return nil
}

func (c *Cnf) ChToScriptsDir() error {
	if err := os.Chdir(c.FullScriptsPath); err != nil {
		return fmt.Errorf("%w\nfull scripts path: %s", err, c.FullScriptsPath)
	}
	return nil
}