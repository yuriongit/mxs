// Package scr manages script data, creation, validation, and related operations.
package scr

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yuriongit/mxs/internal/cfg"
)

// newScriptFixture creates a config with a populated scripts directory and a
// metadata entry for scriptName, returning the config and the script path.
func newScriptFixture(t *testing.T, scriptName string) (*cfg.Cfg, string) {
	t.Helper()

	config := cfg.NewCfg(t.TempDir())
	if err := os.MkdirAll(config.FullScriptsPath, 0700); err != nil {
		t.Fatalf("failed to create scripts dir: %v", err)
	}
	if err := os.MkdirAll(config.FullCfgPath, 0700); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	scriptPath := filepath.Join(config.FullScriptsPath, scriptName+".sh")
	// #nosec G306 -- test fixture writes a small script file
	if err := os.WriteFile(scriptPath, []byte("#!/bin/sh\n"), 0700); err != nil {
		t.Fatalf("failed to write script fixture: %v", err)
	}

	if err := config.RecordLastExecuted(scriptName, time.Date(2026, time.October, 3, 11, 10, 0, 0, time.UTC)); err != nil {
		t.Fatalf("RecordLastExecuted() failed: %v", err)
	}

	return config, scriptPath
}

func TestRemoveScriptDeletesFileAndMetadata(t *testing.T) {
	config, scriptPath := newScriptFixture(t, "backup")

	if err := RemoveScript(config, "backup", false, strings.NewReader("y\n")); err != nil {
		t.Fatalf("RemoveScript() returned unexpected error: %v", err)
	}

	if _, err := os.Stat(scriptPath); !os.IsNotExist(err) {
		t.Fatalf("script file still exists after removal (stat err = %v)", err)
	}

	history, err := config.LoadLastExecuted()
	if err != nil {
		t.Fatalf("LoadLastExecuted() failed: %v", err)
	}
	if _, ok := history["backup"]; ok {
		t.Fatal("metadata[backup] still present after removal")
	}
}

func TestRemoveScriptKeepsFileAndMetadataWhenAborted(t *testing.T) {
	config, scriptPath := newScriptFixture(t, "backup")

	for _, input := range []string{"n\n", "\n", "maybe\n"} {
		if err := RemoveScript(config, "backup", false, strings.NewReader(input)); err != nil {
			t.Fatalf("RemoveScript() on %q returned unexpected error: %v", input, err)
		}

		if _, err := os.Stat(scriptPath); err != nil {
			t.Fatalf("script file missing after aborted removal of %q: %v", input, err)
		}
	}

	history, err := config.LoadLastExecuted()
	if err != nil {
		t.Fatalf("LoadLastExecuted() failed: %v", err)
	}
	if _, ok := history["backup"]; !ok {
		t.Fatal("metadata[backup] removed even though user aborted")
	}
}

func TestRemoveScriptEOFAbortsWithoutDeleting(t *testing.T) {
	config, scriptPath := newScriptFixture(t, "backup")

	// Empty reader simulates EOF (e.g. piped input with no confirmation).
	if err := RemoveScript(config, "backup", false, strings.NewReader("")); err != nil {
		t.Fatalf("RemoveScript() with EOF returned unexpected error: %v", err)
	}

	// The file must remain: EOF aborts without deleting.
	if _, statErr := os.Stat(scriptPath); statErr != nil {
		t.Fatalf("script file missing after EOF abort: %v", statErr)
	}
}

// TestRemoveScriptConfirmWithoutTrailingNewline ensures a bare "y" with no
// newline (e.g. `printf 'y' | mxs rm`) still proceeds, since ReadString
// returns the data together with io.EOF in that case.
func TestRemoveScriptConfirmWithoutTrailingNewline(t *testing.T) {
	config, scriptPath := newScriptFixture(t, "backup")

	if err := RemoveScript(config, "backup", false, strings.NewReader("y")); err != nil {
		t.Fatalf("RemoveScript() with bare 'y' returned unexpected error: %v", err)
	}

	if _, err := os.Stat(scriptPath); !os.IsNotExist(err) {
		t.Fatalf("script file still exists after bare 'y' confirmation")
	}
}

// TestRemoveScriptConfirmVariations confirms the accepted answers are
// case-insensitive and tolerant of surrounding whitespace.
func TestRemoveScriptConfirmVariations(t *testing.T) {
	for _, input := range []string{"Y\n", "yes\n", "YES\n", "  y  \n", " Yes \n"} {
		config, scriptPath := newScriptFixture(t, "backup")

		if err := RemoveScript(config, "backup", false, strings.NewReader(input)); err != nil {
			t.Fatalf("RemoveScript() with %q returned unexpected error: %v", input, err)
		}

		if _, err := os.Stat(scriptPath); !os.IsNotExist(err) {
			t.Fatalf("script file still exists after confirmation %q", input)
		}
	}
}

func TestRemoveScriptForceSkipsConfirmation(t *testing.T) {
	config, scriptPath := newScriptFixture(t, "backup")

	// No input available; force must bypass the prompt entirely.
	if err := RemoveScript(config, "backup", true, strings.NewReader("")); err != nil {
		t.Fatalf("RemoveScript(force=true) returned unexpected error: %v", err)
	}

	if _, err := os.Stat(scriptPath); !os.IsNotExist(err) {
		t.Fatalf("script file still exists after force removal: %v", err)
	}
}

func TestRemoveScriptMissingReturnsError(t *testing.T) {
	config, _ := newScriptFixture(t, "backup")

	err := RemoveScript(config, "deploy", false, strings.NewReader("y\n"))
	if err == nil {
		t.Fatal("RemoveScript() on missing script returned nil error")
	}
}

func TestRemoveScriptRejectsDirectory(t *testing.T) {
	config := cfg.NewCfg(t.TempDir())
	if err := os.MkdirAll(config.FullScriptsPath, 0700); err != nil {
		t.Fatalf("failed to create scripts dir: %v", err)
	}
	if err := os.MkdirAll(config.FullCfgPath, 0700); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	dirPath := filepath.Join(config.FullScriptsPath, "notscript.sh")
	if err := os.Mkdir(dirPath, 0700); err != nil {
		t.Fatalf("failed to create directory fixture: %v", err)
	}

	err := RemoveScript(config, "notscript", false, strings.NewReader("y\n"))
	if err == nil {
		t.Fatal("RemoveScript() deleted a directory named <name>.sh")
	}
	if _, statErr := os.Stat(dirPath); statErr != nil {
		t.Fatalf("directory was removed: %v", statErr)
	}
}

func TestRemoveScriptInvalidNameRejected(t *testing.T) {
	config, _ := newScriptFixture(t, "backup")

	if err := RemoveScript(config, "../evil", false, strings.NewReader("y\n")); err == nil {
		t.Fatal("RemoveScript() accepted an invalid script name")
	}
}

func TestRemoveScriptMetadataFailureAfterFileRemoval(t *testing.T) {
	config, scriptPath := newScriptFixture(t, "backup")

	// Force DeleteMetadataEntry to fail by replacing the metadata file with a
	// directory, which makes reading the store fail.
	if err := os.Remove(config.MetadataPath()); err != nil {
		t.Fatalf("failed to remove metadata file fixture: %v", err)
	}
	if err := os.Mkdir(config.MetadataPath(), 0700); err != nil {
		t.Fatalf("failed to create metadata dir fixture: %v", err)
	}

	err := RemoveScript(config, "backup", false, strings.NewReader("y\n"))
	if err == nil {
		t.Fatal("RemoveScript() succeeded despite metadata write failure")
	}
	if !strings.Contains(err.Error(), "metadata") {
		t.Fatalf("error %q does not surface the metadata failure", err)
	}

	// The file must still be gone: file-first ordering wins.
	if _, statErr := os.Stat(scriptPath); !os.IsNotExist(statErr) {
		t.Fatalf("script file still exists when metadata update failed: %v", statErr)
	}
}
