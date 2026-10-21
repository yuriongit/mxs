/*
Package cnf provides all the functionality
for MXS' config. Manages the application's
configuration and script directory.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cnf

// File represents a user's configuration file for MXS.
// Upon instantiation, a File pointer references an empty struct.
// Use [File.Load] to populate it with the loaded configuration.
type File struct {
	Editor *Editor `json:"editor" validate:"required"`
}

type Editor struct {
	NameAsCommand string `json:"name_as_command" validate:"required"`
	UIType        string `json:"ui_type" validate:"required,oneof=tui gui"`
}
