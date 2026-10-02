package executor

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/yuriongit/xs/internal/cnf"
)

type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

/*
PrepScript verifies the script exists and prepares 
the script future execution.
*/
func (e *Executor) PrepScript(
	cnf *cnf.Cnf,
	scriptName string,
) (scriptCmd *exec.Cmd, err error) {
	// Change to /scripts dir
	if err := cnf.ChToScriptsDir(); err != nil {
		return nil, err
	}
	// Attach the file extension
	fileNameWithExt := e.attachFileExt(scriptName)
	// Check if fileName is a script
	if err := e.isScript(fileNameWithExt); err != nil {
		return nil, err
	}
	// Change file mode if fileNameWithExt is a script
	if err := e.chScriptModToExec(fileNameWithExt); err != nil {
		return nil, fmt.Errorf("failed to change script's file mode: %w", err)
	}

	return exec.Command(fmt.Sprintf("./%s", fileNameWithExt)), nil
}

func (e *Executor) isScript(fileName string) error {
	// Use os.Stat to check if file exists
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

func (e *Executor) chScriptModToExec(fileName string) error {
	return os.Chmod(fileName, 0700)
}