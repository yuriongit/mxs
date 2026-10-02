/* 
Package app wires together the application's 
configuration and execution dependencies.
*/
package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yuriongit/xs/internal/cnf"
	"github.com/yuriongit/xs/internal/executor"
)

func TestNewApp(t *testing.T) {
	homePath := t.TempDir()
	t.Setenv("HOME", homePath)

	app, err := NewApp()
	if err != nil {
		t.Fatalf("NewApp() returned unexpected error: %v", err)
	}

	if app == nil {
		t.Fatal("NewApp() returned nil App")
	}

	if app.HomePath == nil {
		t.Fatal("App.HomePath is nil")
	}

	if *app.HomePath != homePath {
		t.Fatalf(
			"App.HomePath = %q, want %q",
			*app.HomePath,
			homePath,
		)
	}

	if app.Cnf == nil {
		t.Fatal("App.Cnf is nil")
	}

	if app.Executor == nil {
		t.Fatal("App.Executor is nil")
	}

	expectedCnfPath := filepath.Join(homePath, ".xs")
	if app.Cnf.FullCnfPath != expectedCnfPath {
		t.Fatalf(
			"App.Cnf.FullCnfPath = %q, want %q",
			app.Cnf.FullCnfPath,
			expectedCnfPath,
		)
	}

	expectedScriptsPath := filepath.Join(homePath, ".xs", "scripts")
	if app.Cnf.FullScriptsPath != expectedScriptsPath {
		t.Fatalf(
			"App.Cnf.FullScriptsPath = %q, want %q",
			app.Cnf.FullScriptsPath,
			expectedScriptsPath,
		)
	}
}

func TestNewAppInitializesExpectedDependencies(t *testing.T) {
	homePath := t.TempDir()
	t.Setenv("HOME", homePath)

	app, err := NewApp()
	if err != nil {
		t.Fatalf("NewApp() returned unexpected error: %v", err)
	}

	expectedCnf := cnf.NewCnf(homePath)
	if *app.HomePath != homePath {
		t.Fatalf("HomePath = %q, want %q", *app.HomePath, homePath)
	}

	if app.Cnf.Dir != expectedCnf.Dir {
		t.Fatalf("Cnf.Dir = %q, want %q", app.Cnf.Dir, expectedCnf.Dir)
	}

	if app.Cnf.ScriptsPath != expectedCnf.ScriptsPath {
		t.Fatalf(
			"Cnf.ScriptsPath = %q, want %q",
			app.Cnf.ScriptsPath,
			expectedCnf.ScriptsPath,
		)
	}

	if app.Cnf.FullCnfPath != expectedCnf.FullCnfPath {
		t.Fatalf(
			"Cnf.FullCnfPath = %q, want %q",
			app.Cnf.FullCnfPath,
			expectedCnf.FullCnfPath,
		)
	}

	if app.Cnf.FullScriptsPath != expectedCnf.FullScriptsPath {
		t.Fatalf(
			"Cnf.FullScriptsPath = %q, want %q",
			app.Cnf.FullScriptsPath,
			expectedCnf.FullScriptsPath,
		)
	}

	// Confirm that the dependency is the expected concrete executor type.
	var _ *executor.Executor = app.Executor
}

func TestNewAppDoesNotCreateDirectories(t *testing.T) {
	homePath := t.TempDir()
	t.Setenv("HOME", homePath)

	_, err := NewApp()
	if err != nil {
		t.Fatalf("NewApp() returned unexpected error: %v", err)
	}

	cnfPath := filepath.Join(homePath, ".xs")
	scriptsPath := filepath.Join(cnfPath, "scripts")

	if _, err := os.Stat(cnfPath); !os.IsNotExist(err) {
		t.Fatalf(
			"configuration directory was created by NewApp(); stat error = %v",
			err,
		)
	}

	if _, err := os.Stat(scriptsPath); !os.IsNotExist(err) {
		t.Fatalf(
			"scripts directory was created by NewApp(); stat error = %v",
			err,
		)
	}
}