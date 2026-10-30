// Package scr manages script data, creation, validation, and related operations.
package scr

import (
	"testing"
	"time"

	"github.com/yuriongit/mxs/internal/cfg"
)

func TestLastExecPrefix(t *testing.T) {
	executedAt := time.Date(2026, time.October, 30, 11, 10, 0, 0, time.Local)
	history := cfg.LastExecuted{"backup": executedAt}

	prefix, _ := lastExecPrefix("backup", history)
	if want := "l-exec: 11:10am on Oct/30/26 - "; prefix != want {
		t.Fatalf("prefix = %q, want %q", prefix, want)
	}

	prefix, _ = lastExecPrefix("missing", history)
	if want := "l-exec: never executed (n/a) - "; prefix != want {
		t.Fatalf("prefix = %q, want %q", prefix, want)
	}
}
