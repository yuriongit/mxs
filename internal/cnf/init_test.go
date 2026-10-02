/*
Package cnf provides all the functionality
for Xs's config. Manages the application's
configuration and script directory.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cnf

import (
	"os"
	"testing"
)

func TestCnfInitCnfDirCreatesDirectories(t *testing.T) {
	homeDir := t.TempDir()
	config := NewCnf(homeDir)

	if err := config.InitCnfDir(); err != nil {
		t.Fatalf("InitCnfDir() returned unexpected error: %v", err)
	}

	assertDirectory(t, config.FullCnfPath)
	assertDirectory(t, config.FullScriptsPath)
}

func TestCnfInitCnfDirWhenDirectoriesAlreadyExist(t *testing.T) {
	homeDir := t.TempDir()
	config := NewCnf(homeDir)

	if err := os.MkdirAll(config.FullScriptsPath, 0700); err != nil {
		t.Fatalf("failed to create test directories: %v", err)
	}

	if err := config.InitCnfDir(); err != nil {
		t.Fatalf("InitCnfDir() returned unexpected error: %v", err)
	}

	assertDirectory(t, config.FullCnfPath)
	assertDirectory(t, config.FullScriptsPath)
}

func TestCnfInitCnfDirCreatesMissingScriptsDirectory(t *testing.T) {
	homeDir := t.TempDir()
	config := NewCnf(homeDir)

	if err := os.Mkdir(config.FullCnfPath, 0700); err != nil {
		t.Fatalf("failed to create base configuration directory: %v", err)
	}

	if err := config.InitCnfDir(); err != nil {
		t.Fatalf("InitCnfDir() returned unexpected error: %v", err)
	}

	assertDirectory(t, config.FullCnfPath)
	assertDirectory(t, config.FullScriptsPath)
}

func TestCnfInitCnfDirFailsWhenBasePathIsFile(t *testing.T) {
	homeDir := t.TempDir()
	config := NewCnf(homeDir)

	if err := os.WriteFile(
		config.FullCnfPath,
		[]byte("not a directory"),
		0600,
	); err != nil {
		t.Fatalf("failed to create base path file: %v", err)
	}

	err := config.InitCnfDir()
	if err == nil {
		t.Fatal("InitCnfDir() returned nil when base path was a file")
	}

	expected := "~/.xs exists but is not a directory"
	if err.Error() != expected {
		t.Fatalf("error = %q, want %q", err.Error(), expected)
	}
}

func TestCnfInitCnfDirFailsWhenScriptsPathIsFile(t *testing.T) {
	homeDir := t.TempDir()
	config := NewCnf(homeDir)

	if err := os.Mkdir(config.FullCnfPath, 0700); err != nil {
		t.Fatalf("failed to create base configuration directory: %v", err)
	}

	if err := os.WriteFile(
		config.FullScriptsPath,
		[]byte("not a directory"),
		0600,
	); err != nil {
		t.Fatalf("failed to create scripts path file: %v", err)
	}

	err := config.InitCnfDir()
	if err == nil {
		t.Fatal("InitCnfDir() returned nil when scripts path was a file")
	}

	// The implementation builds this display name from c.ScriptsPath,
	// which is "scripts", so the current expected error is ~/scripts.
	expected := "~/scripts exists but is not a directory"
	if err.Error() != expected {
		t.Fatalf("error = %q, want %q", err.Error(), expected)
	}
}

func TestCnfInitCnfDirIsIdempotent(t *testing.T) {
	homeDir := t.TempDir()
	config := NewCnf(homeDir)

	if err := config.InitCnfDir(); err != nil {
		t.Fatalf("first InitCnfDir() call failed: %v", err)
	}

	if err := config.InitCnfDir(); err != nil {
		t.Fatalf("second InitCnfDir() call failed: %v", err)
	}

	assertDirectory(t, config.FullCnfPath)
	assertDirectory(t, config.FullScriptsPath)
}

func assertDirectory(t *testing.T, path string) {
	t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat %q: %v", path, err)
	}

	if !info.IsDir() {
		t.Fatalf("%q exists but is not a directory", path)
	}
}
