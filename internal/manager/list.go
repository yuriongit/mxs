// Package manager provides managerial (CRUD) operations
// for scripts.
package manager

import (
	"os"
	"strings"

	lg "github.com/charmbracelet/lipgloss"
	"github.com/yuriongit/mxs/internal/cnf"
)

func ListAllScripts(cnf *cnf.Cnf) error {
	scriptsDir, err := buildScriptsDirContents(cnf)
	if err != nil {
		return err
	}

	renderScriptsDir(scriptsDir)

	return nil
}

func buildScriptsDirContents(cnf *cnf.Cnf) (string, error) {
	strBuilder := strings.Builder{}

	raw, err := os.ReadDir(cnf.FullScriptsPath)
	if err != nil {
		return "", err
	}

	for idx, val := range raw {
		strBuilder.WriteString(val.Name())
		if (idx + 1) != len(raw) {
			strBuilder.WriteString("\n")
		}
	}

	scriptsDirName := lg.NewStyle().
		PaddingLeft(2).
		Render(
			"~/" + cnf.ScriptsPath + "/",
		)
	scripts := lg.NewStyle().
		PaddingLeft(4).
		Render(
			strBuilder.String(),
		)

	scriptsDirContent := scriptsDirName + "\n" + scripts

	return scriptsDirContent, nil
}
