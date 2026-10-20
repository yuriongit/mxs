/*
Package cnf provides all the functionality
for MXS' config. Manages the application's
configuration and script directory.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cnf

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateNewScript creates a new script file with executable permissions
// and populates it with the configured script template.
// It returns the script's full path or an error if creation fails.
func (c *Cnf) CreateNewScript(scriptName *string) (scriptPath string, err error) {
	if err := c.hasCnfDirectory("scripts"); err != nil {
		return "", err
	}

	file := *scriptName + ".sh"
	fullFilePath := filepath.Join(c.FullScriptsPath, file)

	// Create the file with executable permissions (0700) without overwriting an existing one
	// #nosec G304 G302 -- Create's an executable script with a validated script file path
	f, err := os.OpenFile(fullFilePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0700)
	if err != nil {
		return "", err
	}

	if _, err := f.WriteString(c.fillScriptTemplate(*scriptName)); err != nil {
		return "", err
	}

	// Close the file handle to prevent resource leaks
	return filepath.Join(c.FullScriptsPath, file), f.Close()
}

// fillScriptTemplate populates the configured script template with the given script name.
func (c *Cnf) fillScriptTemplate(scriptName string) string {
	return fmt.Sprintf(
		c.ScriptTemplate, scriptName)
}
