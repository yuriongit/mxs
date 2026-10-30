/*
Package cfg provides all the functionality
for MXS' config. Manages the application's
configuration and script directory.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cfg

import (
	"os"
	"testing"
	"time"
)

func TestLoadLastExecutedMissingFile(t *testing.T) {
	config := NewCfg(t.TempDir())

	history, err := config.LoadLastExecuted()
	if err != nil {
		t.Fatalf("LoadLastExecuted() returned unexpected error: %v", err)
	}

	if history == nil {
		t.Fatal("LoadLastExecuted() returned a nil map")
	}

	if len(history) != 0 {
		t.Fatalf("LoadLastExecuted() returned %d entries, want 0", len(history))
	}
}

func TestRecordLastExecutedRoundTrip(t *testing.T) {
	config := NewCfg(t.TempDir())
	if err := os.MkdirAll(config.FullCfgPath, 0700); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	first := time.Date(2026, time.October, 30, 11, 10, 0, 0, time.UTC)
	second := first.Add(2 * time.Hour)

	if err := config.RecordLastExecuted("backup", first); err != nil {
		t.Fatalf("RecordLastExecuted() returned unexpected error: %v", err)
	}

	if err := config.RecordLastExecuted("deploy", second); err != nil {
		t.Fatalf("RecordLastExecuted() returned unexpected error: %v", err)
	}

	history, err := config.LoadLastExecuted()
	if err != nil {
		t.Fatalf("LoadLastExecuted() returned unexpected error: %v", err)
	}

	if got := history["backup"]; !got.Equal(first) {
		t.Fatalf("history[backup] = %v, want %v", got, first)
	}

	if got := history["deploy"]; !got.Equal(second) {
		t.Fatalf("history[deploy] = %v, want %v", got, second)
	}
}

func TestRecordLastExecutedOverwritesEntry(t *testing.T) {
	config := NewCfg(t.TempDir())
	if err := os.MkdirAll(config.FullCfgPath, 0700); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	first := time.Date(2026, time.October, 30, 11, 10, 0, 0, time.UTC)
	second := first.Add(time.Hour)

	if err := config.RecordLastExecuted("backup", first); err != nil {
		t.Fatalf("RecordLastExecuted() returned unexpected error: %v", err)
	}

	if err := config.RecordLastExecuted("backup", second); err != nil {
		t.Fatalf("RecordLastExecuted() returned unexpected error: %v", err)
	}

	history, err := config.LoadLastExecuted()
	if err != nil {
		t.Fatalf("LoadLastExecuted() returned unexpected error: %v", err)
	}

	if got := history["backup"]; !got.Equal(second) {
		t.Fatalf("history[backup] = %v, want %v", got, second)
	}
}
