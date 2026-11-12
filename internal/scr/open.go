// Package scr manages script data, creation, validation, and related operations.
package scr

import (
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/yuriongit/mxs/internal/cfg"
)

// OpenScript validates the script name, creates the script file,
// opens it in the configured editor, and outputs a success message.
func OpenScript(validator *validator.Validate, cfg *cfg.Cfg, scriptName string) error {
	// Move all config related checks to Pre-Run functions.
	if err := cfg.HasCfgDir("scripts"); err != nil {
		return err
	}
	if err := validateScriptName(scriptName); err != nil {
		return err
	}

	fileName := scriptName + ".sh"
	scriptPath := filepath.Join(cfg.FullScriptsPath, fileName)
	_, err := os.Stat(scriptPath)
	if err != nil {
		return err
	}

	msg, err := OpenEditor("open", validator, cfg, scriptName, scriptPath)
	if err != nil {
		return err
	}

	outputScriptUpdateSuccess(msg)

	return nil
}
