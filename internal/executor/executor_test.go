package executor

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/yuriongit/xs/internal/cnf"
)

func TestExecutorPrepScript(t *testing.T) {
	t.Run("successfully prepares script", func(t *testing.T) {
		homeDir := t.TempDir()
		scriptsDir := filepath.Join(homeDir, ".xs", "scripts")

		if err := os.MkdirAll(scriptsDir, 0700); err != nil {
			t.Fatalf("failed to create scripts directory: %v", err)
		}

		scriptPath := filepath.Join(scriptsDir, "backup.sh")
		scriptContents := []byte("#!/bin/sh\nprintf 'backup complete'\n")

		// #nosec G306 -- enables script execution
		if err := os.WriteFile(scriptPath, scriptContents, 0700); err != nil {
			t.Fatalf("failed to create script: %v", err)
		}

		originalDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get current working directory: %v", err)
		}

		config := cnf.NewCnf(homeDir)
		executor := NewExecutor()

		command, err := executor.PrepScript(config, "backup")
		if err != nil {
			t.Fatalf("PrepScript() returned unexpected error: %v", err)
		}

		if command == nil {
			t.Fatal("PrepScript() returned a nil command")
		}

		expectedScriptPath := filepath.Join(scriptsDir, "backup.sh")

		if command.Path != expectedScriptPath {
			t.Fatalf(
				"command.Path = %q, want %q",
				command.Path,
				expectedScriptPath,
			)
		}

		if command.Dir != scriptsDir {
			t.Fatalf(
				"command.Dir = %q, want %q",
				command.Dir,
				scriptsDir,
			)
		}

		expectedArgs := []string{expectedScriptPath}
		if len(command.Args) != len(expectedArgs) ||
			command.Args[0] != expectedArgs[0] {
			t.Fatalf(
				"command.Args = %#v, want %#v",
				command.Args,
				expectedArgs,
			)
		}

		currentDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get current working directory: %v", err)
		}

		if currentDir != originalDir {
			t.Fatalf(
				"working directory changed from %q to %q",
				originalDir,
				currentDir,
			)
		}

		if runtime.GOOS != "windows" {
			info, err := os.Stat(scriptPath)
			if err != nil {
				t.Fatalf("failed to stat prepared script: %v", err)
			}

			if actualMode := info.Mode().Perm(); actualMode != 0700 {
				t.Fatalf(
					"prepared script permissions = %o, want 0700",
					actualMode,
				)
			}
		}

		output, err := command.Output()
		if err != nil {
			t.Fatalf("returned command failed to execute: %v", err)
		}

		if string(output) != "backup complete" {
			t.Fatalf(
				"command output = %q, want %q",
				string(output),
				"backup complete",
			)
		}
	})

	t.Run("returns an error when scripts directory does not exist", func(t *testing.T) {
		homeDir := t.TempDir()
		config := cnf.NewCnf(homeDir)
		executor := NewExecutor()

		command, err := executor.PrepScript(config, "backup")

		if err == nil {
			t.Fatal("PrepScript() returned nil error")
		}

		if command != nil {
			t.Fatal("PrepScript() returned a command after an error")
		}

		if !strings.Contains(err.Error(), config.FullScriptsPath) {
			t.Fatalf(
				"error = %q, want it to contain %q",
				err.Error(),
				config.FullScriptsPath,
			)
		}
	})

	t.Run("returns an error when script does not exist", func(t *testing.T) {
		homeDir := t.TempDir()
		scriptsDir := filepath.Join(homeDir, ".xs", "scripts")

		if err := os.MkdirAll(scriptsDir, 0700); err != nil {
			t.Fatalf("failed to create scripts directory: %v", err)
		}

		config := cnf.NewCnf(homeDir)
		executor := NewExecutor()

		command, err := executor.PrepScript(config, "missing")

		if err == nil {
			t.Fatal("PrepScript() returned nil error")
		}

		if command != nil {
			t.Fatal("PrepScript() returned a command for a missing script")
		}

		if !os.IsNotExist(err) {
			t.Fatalf(
				"error = %v, want an os.IsNotExist error",
				err,
			)
		}
	})
}
