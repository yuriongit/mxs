// Package scr manages script data, creation, validation, and related operations.
package scr

import (
	"errors"
	"io/fs"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yuriongit/mxs/internal/cfg"
	"github.com/yuriongit/mxs/internal/ui"
)

// ListAllScripts lists out all scripts in the configuration directory.
func ListAllScripts(cfg *cfg.Cfg) error {
	scriptsDir, err := buildScriptsDirContents(cfg)
	if err != nil {
		return err
	}

	renderListAllScriptsView(scriptsDir)

	return nil
}

func buildScriptsDirContents(cfg *cfg.Cfg) (string, error) {
	strBuilder := strings.Builder{}

	raw, err := os.ReadDir(cfg.FullScriptsPath)
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

		// Remove .sh extension
		name := strings.TrimSuffix(val.Name(), ".sh")
		strBuilder.WriteString(name)
	}

	return ui.BaseBox.BorderForeground(lipgloss.Color(ui.Blue)).Foreground(lipgloss.Color(ui.Blue)).Render(strBuilder.String()), nil
}
