package manager

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/yuriongit/mxs/internal/cnf"
	"github.com/yuriongit/mxs/internal/ui"
)

// AddNewScript validates the script name, creates the script file,
// opens it in the configured editor, and outputs a success message.
func AddNewScript(cnf *cnf.Cnf, scriptName *string) error {
	if err := validateScriptName(scriptName); err != nil {
		return err
	}
	scriptPath, err := cnf.CreateNewScript(scriptName)
	if err != nil {
		return err
	}
	msg, err := OpenEditorForScript(scriptName, scriptPath)
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
func OpenEditorForScript(name *string, scriptPath string) (msg string, err error) {
	editorEnvName := "EDITOR"

	editor, ok := os.LookupEnv(editorEnvName)
	if !ok {
		fmt.Println()
		return "", fmt.Errorf("failed to retrieve %s", editorEnvName)
	}

	// #nosec G204 -- opens editor with $EDITOR variable, which does not derive from user input
	cmd := exec.Command(editor, scriptPath)
	if err := cmd.Start(); err != nil {
		return "", err
	}

	msg, err = ui.CreateSuccessMessage(fmt.Sprintf("created %s and opened %s for editing", ui.ScriptName.Render(*name), ui.ScriptDescAndOutroMsgStyle.Render(editor)))
	if err != nil {
		return "", err
	}

	return msg, nil
}

// TODO: Create success message helper
// TODO: Lowercase all messages
