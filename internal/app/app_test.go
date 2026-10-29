/*
Package app wires together the application's
configuration and execution dependencies.
*/
package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yuriongit/mxs/internal/cfg"
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

	if app.Cfg == nil {
		t.Fatal("App.Cfg is nil")
	}

	if app.Executor == nil {
		t.Fatal("App.Executor is nil")
	}

	expectedCfgPath := filepath.Join(homePath, ".mxs")
	if app.Cfg.FullCfgPath != expectedCfgPath {
		t.Fatalf(
			"App.Cfg.FullCfgPath = %q, want %q",
			app.Cfg.FullCfgPath,
			expectedCfgPath,
		)
	}

	expectedScriptsPath := filepath.Join(homePath, ".mxs", "scripts")
	if app.Cfg.FullScriptsPath != expectedScriptsPath {
		t.Fatalf(
			"App.Cfg.FullScriptsPath = %q, want %q",
			app.Cfg.FullScriptsPath,
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

	expectedCfg := cfg.NewCfg(homePath)
	if *app.HomePath != homePath {
		t.Fatalf("HomePath = %q, want %q", *app.HomePath, homePath)
	}

	if app.Cfg.Dir != expectedCfg.Dir {
		t.Fatalf("Cfg.Dir = %q, want %q", app.Cfg.Dir, expectedCfg.Dir)
	}

	if app.Cfg.ScriptsPath != expectedCfg.ScriptsPath {
		t.Fatalf(
			"Cfg.ScriptsPath = %q, want %q",
			app.Cfg.ScriptsPath,
			expectedCfg.ScriptsPath,
		)
	}

	if app.Cfg.FullCfgPath != expectedCfg.FullCfgPath {
		t.Fatalf(
			"Cfg.FullCfgPath = %q, want %q",
			app.Cfg.FullCfgPath,
			expectedCfg.FullCfgPath,
		)
	}

	if app.Cfg.FullScriptsPath != expectedCfg.FullScriptsPath {
		t.Fatalf(
			"Cfg.FullScriptsPath = %q, want %q",
			app.Cfg.FullScriptsPath,
			expectedCfg.FullScriptsPath,
		)
	}

	// Confirm that the dependency is the expected concrete executor type.
	var _ = app.Executor
}

func TestNewAppDoesNotCreateDirectories(t *testing.T) {
	homePath := t.TempDir()
	t.Setenv("HOME", homePath)

	_, err := NewApp()
	if err != nil {
		t.Fatalf("NewApp() returned unexpected error: %v", err)
	}

	cfgPath := filepath.Join(homePath, ".mxs")
	scriptsPath := filepath.Join(cfgPath, "scripts")

	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
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
