// Package manager provides managerial (CRUD) operations
// for scripts.
package manager

import (
	"errors"
	"io/fs"
	"os"
	"strings"

	"github.com/yuriongit/mxs/internal/cnf"
	"github.com/yuriongit/mxs/internal/ui"
)

// ListAllScripts lists out all scripts in the configuration directory.
func ListAllScripts(cnf *cnf.Cnf) error {
	scriptsDir, err := buildScriptsDirContents(cnf)
	if err != nil {
		return err
	}

	renderListAllScriptsView(cnf.ScriptsPath, scriptsDir)

	return nil
}

func buildScriptsDirContents(cnf *cnf.Cnf) (string, error) {
	strBuilder := strings.Builder{}

	raw, err := os.ReadDir(cnf.FullScriptsPath)
	if err != nil {
		return "", err
	} else if len(raw) == 0 {
		return "", errors.New("no existing scripts")
	}

	// Collect only .sh files
	var scriptFiles []fs.DirEntry
	for _, val := range raw {
		if strings.HasSuffix(val.Name(), ".sh") {
			scriptFiles = append(scriptFiles, val)
		}
	}

	if len(scriptFiles) == 0 {
		return "", errors.New("no .sh scripts found")
	}

	// Render with correct symbols
	for idx, val := range scriptFiles {
		if idx > 0 {
			strBuilder.WriteString("\n")
		}

		if idx == len(scriptFiles)-1 {
			strBuilder.WriteString(ui.ScriptEntryStyle.Render(ui.DirEndingEntrySymbol))
		} else {
			strBuilder.WriteString(ui.ScriptEntryStyle.Render(ui.DirEntrySymbol))
		}

		// Remove .sh extension
		name := strings.TrimSuffix(val.Name(), ".sh")
		strBuilder.WriteString(ui.ScriptEntryStyle.Render(name))
	}

	return strBuilder.String(), nil
}
