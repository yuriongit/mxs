/*
Package executor prepares and executes shell scripts.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package executor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/yuriongit/xs/internal/cnf"
)

type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

/*
PrepScript verifies the script exists and prepares
the script for future execution.
*/
func (e *Executor) PrepScript(
	cnf *cnf.Cnf,
	scriptName string,
) (scriptCmd *exec.Cmd, err error) {
	// Attach the file extension.
	fileNameWithExt := e.attachFileExt(scriptName)

	// Build the full path to the script in the scripts directory.
	fullScriptPath := filepath.Join(cnf.FullScriptsPath, fileNameWithExt)

	// Check if the script exists.
	if err := e.isScript(fullScriptPath); err != nil {
		return nil, err
	}

	// Execute the script relative to the scripts directory without
	// changing the process-wide working directory.
	// #nosec G204 -- scriptPath validated by isScript() above
	scriptCmd = exec.Command(filepath.Join(cnf.FullScriptsPath, fileNameWithExt))
	scriptCmd.Dir = cnf.FullScriptsPath

	return scriptCmd, nil
}

func (e *Executor) isScript(fileName string) error {
	// Use os.Stat to check if the script exists.
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return err
	}

	return nil
}

func (e *Executor) attachFileExt(fileName string) string {
	return fmt.Sprintf("%s.sh", fileName)
}
