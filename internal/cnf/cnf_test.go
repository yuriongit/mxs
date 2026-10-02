package cnf

import (
	"path/filepath"
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
