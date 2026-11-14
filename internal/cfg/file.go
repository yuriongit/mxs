/*
Package cfg provides all the functionality
for MXS' config. Manages the application's
configuration and script directory.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cfg

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/yuriongit/mxs/internal/ui"
)

// File represents a user's configuration file for MXS.
// Upon instantiation, a File pointer references an empty struct.
// Use [File.Load] to populate it with the loaded configuration.
//
// Handles configuration decoding and validation.
type File struct {
	Editor *Editor `json:"editor" validate:"required"`
}

// Editor defines the needed information for editor usage.
//
// Needed for the add and upcoming upd command.
type Editor struct {
	NameAsCommand string `json:"name_as_command" validate:"required"`
	UIType        string `json:"ui_type" validate:"required,oneof=tui gui"`
}

// Load reads and validates the configuration file.
func (f *File) Load(cfg *Cfg, validate *validator.Validate) error {
	if err := f.unmarshal(cfg); err != nil {
		return err
	}

	if err := validate.Struct(f); err != nil {
		return fmt.Errorf("configuration errors:\n%w", err)
	}

	return nil
}

func (f *File) unmarshal(cfg *Cfg) error {
	cfgFilePath := filepath.Join(cfg.FullCfgPath, cfg.FileName)
	//#nosec: G304 -- safe from user input; variable uses predetermined values
	rawJSON, err := os.ReadFile(cfgFilePath)
	if err != nil {
		return fmt.Errorf("reading configuration file: %w", err)
	}

	if err := json.Unmarshal(rawJSON, f); err != nil {
		return fmt.Errorf("decoding configuration file: %w", err)
	}

	return nil
}

// Open opens the configuration file in the user's configured editor.
//
// Revisit to refactor [File.Open] and [Scr.OpenEditor] into one generic
// (create editor.go file in package cfg).
func (f *File) Open(cfg *Cfg, validator *validator.Validate) (msg string, err error) {
	if err := f.Load(cfg, validator); err != nil {
		return "", err
	}
	if cfg.File == nil {
		return "", fmt.Errorf("cfg file is a nil pointer")
	}

	editor := cfg.File.Editor

	// #nosec G204 -- comes from the user's config file for their own editor command
	// (additionally, MXS will only run on their hardware)
	//
	// Additionally, cfg.FilePath is already predetermined in the cfg package (trusted & safe).
	cmd := exec.Command(editor.NameAsCommand, cfg.FilePath)

	switch editor.UIType {
	case "gui":
		if err := cmd.Start(); err != nil {
			return "", err
		}
		return ui.CreateSuccMsg(fmt.Sprintf(
			"opened %s with %s",
			ui.ScriptReference.Render("config"),
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
			"opened %s with %s",
			ui.ScriptReference.Render("config"),
			ui.Primary.Render(editor.NameAsCommand)),
		), nil
	default:
		return
	}
}
