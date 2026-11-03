package executor

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestLineWriterBuffersPartialWrites(t *testing.T) {
	sub := make(chan tea.Msg, 16)
	writer := &lineWriter{sub: sub}

	if _, err := writer.Write([]byte("partial")); err != nil {
		t.Fatalf("Write() returned unexpected error: %v", err)
	}

	select {
	case msg := <-sub:
		t.Fatalf("Write() forwarded incomplete line: %q", string(msg.(outputLineMsg)))
	default:
	}

	if _, err := writer.Write([]byte("line1\n")); err != nil {
		t.Fatalf("Write() returned unexpected error: %v", err)
	}

	select {
	case msg := <-sub:
		if got := string(msg.(outputLineMsg)); got != "partialline1" {
			t.Fatalf("forwarded line = %q, want %q", got, "partialline1")
		}
	default:
		t.Fatal("Write() did not forward completed line")
	}

	if _, err := writer.Write([]byte("line2\n")); err != nil {
		t.Fatalf("Write() returned unexpected error: %v", err)
	}

	select {
	case msg := <-sub:
		if got := string(msg.(outputLineMsg)); got != "line2" {
			t.Fatalf("forwarded line = %q, want %q", got, "line2")
		}
	default:
		t.Fatal("Write() did not forward completed line")
	}
}

func TestRunScriptCmdStreamsOutputLinesInOrder(t *testing.T) {
	scriptsDir := t.TempDir()
	scriptPath := filepath.Join(scriptsDir, "mixed.sh")
	scriptContents := []byte("#!/bin/bash\necho out1\necho err1 >&2\necho out2\n")

	// #nosec G306 -- enables script execution
	if err := os.WriteFile(scriptPath, scriptContents, 0700); err != nil {
		t.Fatalf("failed to create script: %v", err)
	}

	sub := make(chan tea.Msg, 16)
	finished := make(chan error, 1)

	cmd := runScriptCmd(scriptPath, nil, sub)
	go func() {
		msg := cmd()
		finished <- msg.(executionFinishedMsg).err
	}()

	var lines []string
	for {
		select {
		case msg, ok := <-sub:
			if !ok {
				t.Fatal("output channel closed before execution finished")
			}
			if line, ok := msg.(outputLineMsg); ok {
				lines = append(lines, string(line))
			}
		case err := <-finished:
			if err != nil {
				t.Fatalf("script execution failed: %v", err)
			}
			expected := []string{"out1", "err1", "out2"}
			if len(lines) != len(expected) {
				t.Fatalf("streamed lines = %v, want %v", lines, expected)
			}
			for i, line := range lines {
				if line != expected[i] {
					t.Fatalf("streamed line %d = %q, want %q", i, line, expected[i])
				}
			}
			return
		}
	}
}
