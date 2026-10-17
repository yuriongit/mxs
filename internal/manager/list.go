// Package manager provides managerial (CRUD) operations
// for scripts.
package manager

import (
	"errors"
	"os"
	"strings"

	"github.com/yuriongit/mxs/internal/cnf"
	"github.com/yuriongit/mxs/internal/ui"
)

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
		return "", errors.New("No existing scripts")
	}

	for idx, val := range raw {
		switch {
		case !strings.HasSuffix(val.Name(), ".sh"):
			continue
		case (idx + 1) != len(raw):
			strBuilder.WriteString(ui.ScriptEntryStyle.Render(ui.DirEntrySymbol))
			strBuilder.WriteString(ui.ScriptEntryStyle.Render(strings.TrimSuffix(val.Name(), ".sh")))
			strBuilder.WriteString("\n")
		default:
			strBuilder.WriteString(ui.ScriptEntryStyle.Render(ui.DirEndingEntrySymbol))
			strBuilder.WriteString(ui.ScriptEntryStyle.Render(strings.TrimSuffix(val.Name(), ".sh")))
		}
	}

	return strBuilder.String(), nil
}

// TODO: Ignore non .sh files in the scripts directory.
// if !strings.Contains(val.Name(), ".sh") {
// 			continue
// 		} else {
