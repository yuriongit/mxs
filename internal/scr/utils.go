// Package scr manages script data, creation, validation, and related operations.
package scr

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/go-playground/validator/v10"
	"github.com/yuriongit/mxs/internal/cfg"
	"github.com/yuriongit/mxs/internal/ui"
)

// OpenEditor launches the configured editor with the script file.
// It returns a success message or an error if the editor cannot be launched.
func OpenEditor(
	purpose string,
	validator *validator.Validate,
	cfg *cfg.Cfg,
	name string,
	scriptPath string,
) (msg string, err error) {
	if err := cfg.File.Load(cfg, validator); err != nil {
		return "", err
	}
	if cfg.File == nil {
		return "", fmt.Errorf("cfg file is a nil pointer")
	}

	editor := cfg.File.Editor

	// #nosec G204 -- comes from the user's config file for their own editor command
	// (additionally, MXS will only run on their hardware)
	cmd := exec.Command(editor.NameAsCommand, scriptPath)

	var messageTemplate string
	switch purpose {
	case "open":
		messageTemplate = "opened %s with %s"
	case "creation":
		messageTemplate = "created and opened %s with %s"
	default:
		panic(fmt.Sprint("invaild purpose provided for opening user's editor: ", purpose))
	}

	switch editor.UIType {
	case "gui":
		if err := cmd.Start(); err != nil {
			return "", err
		}
		return ui.CreateSuccMsg(fmt.Sprintf(
			messageTemplate,
			ui.ScriptReference.Render(name),
			ui.Primary.Render(editor.NameAsCommand)),
		), nil
	case "tui":
		// Temporarily reassigns parent's standard input, output, and error streams to child process (e.g, hands off to Neovim, then returns back to parent)
		// Bubbletea is not being used here, Cobra is, so this is fine.
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		// When Bubbletea is being used, it provides tea.ExecProcess() which is
		// specifically for launching external terminal apps, it does the following:
		//  1. Suspends Bubble Tea's rendering
		//  2. Releases its terminal input handling
		//  3. Executes the external process
		//  4. Restores the terminal when the process exits

		if err := cmd.Run(); err != nil {
			return "", err
		}
		return ui.CreateSuccMsg(fmt.Sprintf(
			messageTemplate,
			ui.ScriptReference.Render(name),
			ui.Primary.Render(editor.NameAsCommand)),
		), nil
	default:
		return
	}
}
