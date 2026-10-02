/*
Package cnf provides all the functionality
for Xs's config. Manages the application's
configuration and script directory.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cnf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCnfDirCreatesDirectories(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	if err := InitCnfDir(); err != nil {
		t.Fatalf("InitCnfDir() returned unexpected error: %v", err)
	}

	assertDirectory(t, filepath.Join(homeDir, ".xs"))
	assertDirectory(t, filepath.Join(homeDir, ".xs", "scripts"))
}

func TestInitCnfDirWhenDirectoriesAlreadyExist(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	xsDir := filepath.Join(homeDir, ".xs")
	scriptsDir := filepath.Join(xsDir, "scripts")

	if err := os.MkdirAll(scriptsDir, 0700); err != nil {
		t.Fatalf("failed to create test directories: %v", err)
	}

	if err := InitCnfDir(); err != nil {
		t.Fatalf("InitCnfDir() returned unexpected error: %v", err)
	}

	assertDirectory(t, xsDir)
	assertDirectory(t, scriptsDir)
}

func TestInitCnfDirCreatesMissingScriptsDirectory(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	xsDir := filepath.Join(homeDir, ".xs")
	scriptsDir := filepath.Join(xsDir, "scripts")

	if err := os.Mkdir(xsDir, 0700); err != nil {
		t.Fatalf("failed to create .xs directory: %v", err)
	}

	if err := InitCnfDir(); err != nil {
		t.Fatalf("InitCnfDir() returned unexpected error: %v", err)
	}

	assertDirectory(t, xsDir)
	assertDirectory(t, scriptsDir)
}

func TestInitCnfDirFailsWhenXsPathIsFile(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	xsPath := filepath.Join(homeDir, ".xs")
	if err := os.WriteFile(xsPath, []byte("not a directory"), 0600); err != nil {
		t.Fatalf("failed to create .xs file: %v", err)
	}

	err := InitCnfDir()
	if err == nil {
		t.Fatal("InitCnfDir() returned nil error when .xs is a file")
	}

	expected := "~/.xs exists but is not a directory"
	if err.Error() != expected {
		t.Fatalf("error = %q, want %q", err.Error(), expected)
	}
}

func TestInitCnfDirFailsWhenScriptsPathIsFile(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	xsDir := filepath.Join(homeDir, ".xs")
	scriptsPath := filepath.Join(xsDir, "scripts")

	if err := os.Mkdir(xsDir, 0700); err != nil {
		t.Fatalf("failed to create .xs directory: %v", err)
	}

	if err := os.WriteFile(scriptsPath, []byte("not a directory"), 0600); err != nil {
		t.Fatalf("failed to create scripts file: %v", err)
	}

	err := InitCnfDir()
	if err == nil {
		t.Fatal("InitCnfDir() returned nil error when scripts is a file")
	}

	expected := "~/.xs/scripts exists but is not a directory"
	if err.Error() != expected {
		t.Fatalf("error = %q, want %q", err.Error(), expected)
	}
}

func TestInitCnfDirIsIdempotent(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	if err := InitCnfDir(); err != nil {
		t.Fatalf("first InitCnfDir() call failed: %v", err)
	}

	if err := InitCnfDir(); err != nil {
		t.Fatalf("second InitCnfDir() call failed: %v", err)
	}

	assertDirectory(t, filepath.Join(homeDir, ".xs"))
	assertDirectory(t, filepath.Join(homeDir, ".xs", "scripts"))
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
