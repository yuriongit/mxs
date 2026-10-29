// Package scr manages script data, creation, validation, and related operations.
package scr

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuriongit/mxs/internal/cfg"
	"github.com/yuriongit/mxs/internal/ui"
)

// ListAllScripts lists out all scripts in the configuration directory.
func ListAllScripts(cfg *cfg.Cfg, scr *Scr, verbose bool) error {
	scriptsDir, err := buildScriptsDirContents(cfg, scr, verbose)
	if err != nil {
		return err
	}

	renderListAllScriptsView(scriptsDir)

	return nil
}

func buildScriptsDirContents(cfg *cfg.Cfg, scr *Scr, verbose bool) (string, error) {
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

	var rows []ui.ColumnPair
	for _, val := range scriptFiles {
		name := strings.TrimSuffix(val.Name(), ".sh")

		if verbose {
			scriptPath := filepath.Join(cfg.FullScriptsPath, val.Name())
			// #nosec G304 -- safe from user input
			file, fileReadErr := os.ReadFile(scriptPath)
			if fileReadErr != nil {
				return "", fileReadErr
			}

			desc, extractDescErr := scr.ExtractDesc(string(file))
			if extractDescErr != nil {
				desc = ui.Subtle.Render("n/a")
			}

			rows = append(rows, ui.ColumnPair{
				Left:      name,
				Right:     desc,
				LeftStyle: ui.Primary,
			})
		} else {
			rows = append(rows, ui.ColumnPair{
				Left:      name,
				LeftStyle: ui.Primary,
			})
		}
	}

	// 4 defines the column gap space between name and desc
	content := ui.FormatAlignedColumns(rows, 4)

	return ui.BaseBox.
		BorderForeground(ui.Primary.GetForeground()).
		Render(content), nil
}
