// Package scr manages script data, creation, validation, and related operations.
package scr

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	lg "github.com/charmbracelet/lipgloss"
	"github.com/yuriongit/mxs/internal/cfg"
	"github.com/yuriongit/mxs/internal/ui"
)

// lastExecTimeLayout is the display format for a script's last execution time.
const lastExecTimeLayout = "3:04pm on Jan/02/06"

// ListAllScripts lists out all scripts in the configuration directory.
func ListAllScripts(config *cfg.Cfg, scr *Scr, verbose bool) error {
	var history cfg.LastExecuted
	if verbose {
		var err error
		history, err = config.LoadLastExecuted()
		if err != nil {
			return err
		}
	}

	scriptsDir, err := buildScriptsDirContents(config, scr, history, verbose)
	if err != nil {
		return err
	}

	renderListAllScriptsView(scriptsDir)

	return nil
}

// lastExecPrefix builds the last-executed tag and its style for a script.
//
// Scripts with no recorded execution use [ui.Subtle] for the tag.
func lastExecPrefix(name string, history cfg.LastExecuted) (string, lg.Style) {
	if executedAt, ok := history[name]; ok && !executedAt.IsZero() {
		return fmt.Sprintf("l-exec: %s - ", executedAt.Format(lastExecTimeLayout)), lg.NewStyle()
	}

	return "l-exec: never executed (n/a) - ", ui.Subtle
}

func buildScriptsDirContents(
	cfg *cfg.Cfg,
	scr *Scr,
	history cfg.LastExecuted,
	verbose bool,
) (string, error) {
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

			prefix, prefixStyle := lastExecPrefix(name, history)
			rows = append(rows, ui.ColumnPair{
				Left:            name,
				Right:           desc,
				LeftStyle:       ui.Primary,
				LeftPrefix:      prefix,
				LeftPrefixStyle: prefixStyle,
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
