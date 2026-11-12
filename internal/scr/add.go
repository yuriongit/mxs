// Package scr manages script data, creation, validation, and related operations.
package scr

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/yuriongit/mxs/internal/cfg"
	"github.com/yuriongit/mxs/internal/ui"
)

// AddNewScript validates the script name, creates the script file,
// opens it in the configured editor, and outputs a success message.
func AddNewScript(validator *validator.Validate, cfg *cfg.Cfg, scr *Scr, scriptName *string) error {
	if err := cfg.HasCfgDir("scripts"); err != nil {
		return err
	}

	if err := validateScriptName(scriptName); err != nil {
		return err
	}
	scriptPath, err := scr.Create(cfg, scriptName)
	if err != nil {
		return err
	}
	msg, err := OpenEditorForScript(validator, cfg, scriptName, scriptPath)
	if err != nil {
		return err
	}

	outputScriptCreationSuccess(msg)

	return nil
}

// validateScriptName checks that the script name meets the required naming rules.
func validateScriptName(name *string) error {
	// Check length
	if len(*name) == 0 {
		return fmt.Errorf("script name required")
	}
	if len(*name) > 63 {
		return fmt.Errorf("script name must be 63 characters or less (got %d)", len(*name))
	}

	// Check first character (must be lowercase letter, +, -)
	firstChar := rune((*name)[0])
	if firstChar != '+' && firstChar != '-' && (firstChar < 'a' || firstChar > 'z') {
		return fmt.Errorf("script name must start with a lowercase letter or '+' (got '%c')", firstChar)
	}

	// Check all characters
	for i, ch := range *name {
		isLowercase := ch >= 'a' && ch <= 'z'
		isDigit := ch >= '0' && ch <= '9'
		isHyphen := ch == '-'
		isUnderscore := ch == '_'
		isPlus := ch == '+'

		if !isLowercase && !isDigit && !isHyphen && !isUnderscore && !isPlus {
			return fmt.Errorf("script name contains invalid character '%c' (position %d): only lowercase letters, digits, hyphens, underscores, and plus signs are allowed", ch, i+1)
		}
	}

	return nil
}

// OpenEditorForScript launches the configured editor with the script file.
// It returns a success message or an error if the editor cannot be launched.
func OpenEditorForScript(validator *validator.Validate, cfg *cfg.Cfg, name *string, scriptPath string) (msg string, err error) {
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

	switch editor.UIType {
	case "gui":
		if err := cmd.Start(); err != nil {
			return "", err
		}
		return ui.CreateSuccMsg(fmt.Sprintf(
			"created %s and opened %s for editing",
			ui.ScriptReference.Render(*name),
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
			"created %s and opened %s for editing",
			ui.ScriptReference.Render(*name),
			ui.Primary.Render(editor.NameAsCommand)),
		), nil
	default:
		return
	}
}

// Create creates a new script with executable permissions
// (0700 mode). The script is already written with a small
// template that includes the Bash shebang, script header,
// and a statement that echos the script name.
//
// Create returns the script's full path or an error if
// creation fails.
func (s *Scr) Create(cfg *cfg.Cfg, scriptName *string) (scriptPath string, err error) {
	fileName := *scriptName + ".sh"
	fullFilePath := filepath.Join(cfg.FullScriptsPath, fileName)

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
	return filepath.Join(cfg.FullScriptsPath, fileName), file.Close()
}

// writeTemplate fills out the script template with the script name.
func (s *Scr) writeTemplate(scriptName string) string {
	return fmt.Sprintf(
		s.Template, TemplateDesc, scriptName,
	)
}
