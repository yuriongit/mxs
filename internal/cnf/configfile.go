/*
Package cnf provides all the functionality
for MXS' config. Manages the application's
configuration and script directory.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cnf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
)

// File represents a user's configuration file for MXS.
// Upon instantiation, a File pointer references an empty struct.
// Use [File.Load] to populate it with the loaded configuration.
//
// Handles configuration decoding and validation.
type File struct {
	Editor *Editor `json:"editor" validate:"required"`
}

type Editor struct {
	NameAsCommand string `json:"name_as_command" validate:"required"`
	UIType        string `json:"ui_type" validate:"required,oneof=tui gui"`
}

func (f *File) Load(cnf *Cnf, validate *validator.Validate) error {
	if err := f.unmarshal(cnf); err != nil {
		return err
	}

	if err := validate.Struct(f); err != nil {
		return fmt.Errorf("configuration errors:\n%w", err)
	}

	return nil
}

func (f *File) unmarshal(cnf *Cnf) error {
	cnfFilePath := filepath.Join(cnf.FullCnfPath, cnf.FileName)
	//#nosec: G304 -- safe from user input; variable uses predetermined values
	rawJSON, err := os.ReadFile(cnfFilePath)
	if err != nil {
		return fmt.Errorf("reading configuration file: %w", err)
	}

	if err := json.Unmarshal(rawJSON, f); err != nil {
		return fmt.Errorf("decoding configuration file: %w", err)
	}

	return nil
}
