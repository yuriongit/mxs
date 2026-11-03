// Package scr manages script data, creation, validation, and related operations.
package scr

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuriongit/mxs/internal/cfg"
	"github.com/yuriongit/mxs/internal/ui"
)

// RemoveScript deletes the script file for scriptName and removes its metadata
// entry from the store.
//
// The script file is removed first, then its metadata entry, so the file is
// never left behind. If the metadata update fails afterward, the error is
// wrapped to make clear the file was already removed.
//
// When force is false, the caller is prompted to confirm on in; a "y"/"yes"
// (case-insensitive) proceeds, anything else aborts without error.
func RemoveScript(config *cfg.Cfg, scriptName string, force bool, in io.Reader) error {
	if err := validateScriptName(scriptName); err != nil {
		return err
	}

	scriptPath := filepath.Join(config.FullScriptsPath, scriptName+".sh")

	// Reject directories named <name>.sh and report a missing script before
	// attempting any removal.
	info, statErr := os.Stat(scriptPath)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return fmt.Errorf("no script named %q", scriptName)
		}
		return fmt.Errorf("accessing %s: %w", scriptPath, statErr)
	}
	if info.IsDir() {
		return fmt.Errorf("%s exists but is a directory, not a script", scriptPath)
	}

	if !force {
		if !confirmRemoval(scriptName, in) {
			return nil // aborted by the user; not an error
		}
	}

	// #nosec G304 -- path built from the validated script name and config dir
	if err := os.Remove(scriptPath); err != nil {
		return fmt.Errorf("removing %s: %w", scriptPath, err)
	}

	if err := config.DeleteMetadataEntry(scriptName); err != nil {
		return fmt.Errorf(
			"removed %s but failed to update metadata: %w",
			ui.ScrName.Render(scriptName), err,
		)
	}

	outputScriptRemovalSuccess(scriptName)

	return nil
}

// confirmRemoval prompts the user to confirm deletion of scriptName by reading
// a single line from in.
//
// It returns true only on an explicit "y"/"yes" (case-insensitive). EOF, empty
// input, "n", or any malformed input aborts and returns false.
func confirmRemoval(scriptName string, in io.Reader) bool {
	fmt.Fprintf(
		os.Stderr,
		"Are you sure you want to delete %s? [y/N]: ",
		ui.ScrName.Render(scriptName),
	)

	line, err := bufio.NewReader(in).ReadString('\n')
	// ReadString returns the data read together with io.EOF when the input has
	// no trailing newline (e.g. `printf 'y'`). Honor that data instead of
	// treating it as an empty/aborted response; only fail closed when we truly
	// read nothing (empty input or a non-EOF read error).
	if err != nil && err != io.EOF && line == "" {
		return false
	}

	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
		return true
	default:
		return false
	}
}

// outputScriptRemovalSuccess prints the script removal success message.
func outputScriptRemovalSuccess(scriptName string) {
	fmt.Println(ui.CreateSuccMsg(fmt.Sprintf(
		"removed %s",
		ui.ScrName.Render(scriptName),
	)))
}
