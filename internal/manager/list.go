// Package manager provides managerial (CRUD) operations
// for scripts.
package manager

import (
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
	}

	for idx, val := range raw {
		if (idx + 1) != len(raw) {
			strBuilder.WriteString(ui.ScriptEntryStyle.Render(ui.DirEntrySymbol))
			strBuilder.WriteString(ui.ScriptEntryStyle.Render(val.Name()))
			strBuilder.WriteString("\n")
		} else {
			strBuilder.WriteString(ui.ScriptEntryStyle.Render(ui.DirEndingEntrySymbol))
			strBuilder.WriteString(ui.ScriptEntryStyle.Render(val.Name()))
		}
	}

	return strBuilder.String(), nil
}
