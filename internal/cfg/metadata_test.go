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

func TestLoadMetadataMissingFile(t *testing.T) {
	config := NewCfg(t.TempDir())

	metadata, err := config.LoadMetadata()
	if err != nil {
		t.Fatalf("LoadMetadata() returned unexpected error: %v", err)
	}

	if len(metadata) != 0 {
		t.Fatalf("LoadMetadata() returned %d entries, want 0", len(metadata))
	}
}

func TestSaveAndLoadMetadataRoundTrip(t *testing.T) {
	config := NewCfg(t.TempDir())
	if err := os.MkdirAll(config.FullCfgPath, 0700); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	record := time.Date(2026, time.October, 30, 11, 10, 0, 0, time.UTC)
	if err := config.SaveMetadata(Metadata{"backup": record}); err != nil {
		t.Fatalf("SaveMetadata() returned unexpected error: %v", err)
	}

	metadata, err := config.LoadMetadata()
	if err != nil {
		t.Fatalf("LoadMetadata() returned unexpected error: %v", err)
	}

	if got, ok := metadata["backup"].(string); !ok || got == "" {
		t.Fatalf("metadata[backup] = %v (%T), want a non-empty RFC3339 string", metadata["backup"], metadata["backup"])
	}
}

func TestDeleteMetadataEntryRemovesOnlyTarget(t *testing.T) {
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

	if err := config.DeleteMetadataEntry("backup"); err != nil {
		t.Fatalf("DeleteMetadataEntry() returned unexpected error: %v", err)
	}

	metadata, err := config.LoadLastExecuted()
	if err != nil {
		t.Fatalf("LoadLastExecuted() returned unexpected error: %v", err)
	}

	if _, ok := metadata["backup"]; ok {
		t.Fatal("metadata[backup] still present after DeleteMetadataEntry")
	}

	if got, ok := metadata["deploy"]; !ok || !got.Equal(second) {
		t.Fatalf("metadata[deploy] = %v (present=%v), want %v", got, ok, second)
	}
}

func TestDeleteMetadataEntryMissingIsNoError(t *testing.T) {
	config := NewCfg(t.TempDir())
	if err := os.MkdirAll(config.FullCfgPath, 0700); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	if err := config.DeleteMetadataEntry("ghost"); err != nil {
		t.Fatalf("DeleteMetadataEntry() on missing entry returned error: %v", err)
	}

	metadata, err := config.LoadMetadata()
	if err != nil {
		t.Fatalf("LoadMetadata() returned unexpected error: %v", err)
	}
	if len(metadata) != 0 {
		t.Fatalf("LoadMetadata() returned %d entries, want 0", len(metadata))
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
