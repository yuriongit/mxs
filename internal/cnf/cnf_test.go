package cnf

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewCnf(t *testing.T) {
	homePath := filepath.Join("home", "yuri")

	config := NewCnf(homePath)

	if config == nil {
		t.Fatal("NewCnf() returned nil")
	}

	if config.Dir != ".xs" {
		t.Fatalf("Dir = %q, want %q", config.Dir, ".xs")
	}

	if config.ScriptsPath != "scripts" {
		t.Fatalf(
			"ScriptsPath = %q, want %q",
			config.ScriptsPath,
			"scripts",
		)
	}

	expectedCnfPath := filepath.Join(homePath, ".xs")
	if config.FullCnfPath != expectedCnfPath {
		t.Fatalf(
			"FullCnfPath = %q, want %q",
			config.FullCnfPath,
			expectedCnfPath,
		)
	}

	expectedScriptsPath := filepath.Join(homePath, ".xs", "scripts")
	if config.FullScriptsPath != expectedScriptsPath {
		t.Fatalf(
			"FullScriptsPath = %q, want %q",
			config.FullScriptsPath,
			expectedScriptsPath,
		)
	}
}

func TestChToBaseDir(t *testing.T) {
	t.Run("changes to existing base directory", func(t *testing.T) {
		homePath := t.TempDir()
		basePath := filepath.Join(homePath, ".xs")

		if err := os.Mkdir(basePath, 0700); err != nil {
			t.Fatalf("failed to create base directory: %v", err)
		}

		restoreWorkingDirectory(t)

		config := NewCnf(homePath)

		if err := config.ChToBaseDir(); err != nil {
			t.Fatalf("ChToBaseDir() returned unexpected error: %v", err)
		}

		currentPath, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get current working directory: %v", err)
		}

		if currentPath != basePath {
			t.Fatalf(
				"current directory = %q, want %q",
				currentPath,
				basePath,
			)
		}
	})

	t.Run("returns an error when base directory does not exist", func(t *testing.T) {
		homePath := t.TempDir()
		config := NewCnf(homePath)

		originalPath, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get current working directory: %v", err)
		}

		err = config.ChToBaseDir()
		if err == nil {
			t.Fatal("ChToBaseDir() returned nil for a missing directory")
		}

		if !strings.Contains(err.Error(), config.FullCnfPath) {
			t.Fatalf(
				"error = %q, want it to contain %q",
				err.Error(),
				config.FullCnfPath,
			)
		}

		currentPath, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get current working directory: %v", err)
		}

		if currentPath != originalPath {
			t.Fatalf(
				"working directory changed after failed ChToBaseDir(): %q -> %q",
				originalPath,
				currentPath,
			)
		}
	})
}

func TestChToScriptsDir(t *testing.T) {
	t.Run("changes to existing scripts directory", func(t *testing.T) {
		homePath := t.TempDir()
		scriptsPath := filepath.Join(homePath, ".xs", "scripts")

		if err := os.MkdirAll(scriptsPath, 0700); err != nil {
			t.Fatalf("failed to create scripts directory: %v", err)
		}

		restoreWorkingDirectory(t)

		config := NewCnf(homePath)

		if err := config.ChToScriptsDir(); err != nil {
			t.Fatalf("ChToScriptsDir() returned unexpected error: %v", err)
		}

		currentPath, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get current working directory: %v", err)
		}

		if currentPath != scriptsPath {
			t.Fatalf(
				"current directory = %q, want %q",
				currentPath,
				scriptsPath,
			)
		}
	})

	t.Run("returns an error when scripts directory does not exist", func(t *testing.T) {
		homePath := t.TempDir()
		config := NewCnf(homePath)

		originalPath, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get current working directory: %v", err)
		}

		err = config.ChToScriptsDir()
		if err == nil {
			t.Fatal("ChToScriptsDir() returned nil for a missing directory")
		}

		if !strings.Contains(err.Error(), config.FullScriptsPath) {
			t.Fatalf(
				"error = %q, want it to contain %q",
				err.Error(),
				config.FullScriptsPath,
			)
		}

		currentPath, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get current working directory: %v", err)
		}

		if currentPath != originalPath {
			t.Fatalf(
				"working directory changed after failed ChToScriptsDir(): %q -> %q",
				originalPath,
				currentPath,
			)
		}
	})
}

// restoreWorkingDirectory restores the process-wide working directory after
// a test. os.Chdir affects the entire test process, so these tests must not
// be run in parallel.
func restoreWorkingDirectory(t *testing.T) {
	t.Helper()

	originalPath, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(originalPath); err != nil {
			t.Errorf(
				"failed to restore working directory to %q: %v",
				originalPath,
				err,
			)
		}
	})
}
