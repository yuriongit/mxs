package scr

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yuriongit/mxs/internal/cnf"
)

// Create creates a new script with executable permissions
// (0700 mode). The script is already written with a small 
// template that includes the Bash shebang, script header, 
// and a statement that echos the script name.
//
// Create returns the script's full path or an error if 
// creation fails.
func (s *Scr) Create(cnf *cnf.Cnf,scriptName *string) (scriptPath string, err error) {
	fileName := *scriptName + ".sh"
	fullFilePath := filepath.Join(cnf.FullScriptsPath, fileName)

	// Create the file with executable permissions (0700) without overwriting an existing one
	// #nosec G304 G302 -- Create's an executable script with a validated script file path
	file, err := os.OpenFile(fullFilePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0700)
	if err != nil {
		return "", err
	}

	if _, err := file.WriteString(s.writeTemplate(*scriptName)); err != nil {
		return "", err
	}

	// Close the file handle to prevent resource leaks
	return filepath.Join(cnf.FullScriptsPath, fileName), file.Close()
}

// writeTemplate fills out the script template with the script name.
func (s *Scr) writeTemplate(scriptName string) string {
	return fmt.Sprintf(
		s.Template, scriptName,
	)
}
