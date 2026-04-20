package runtime

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
)

func TestTaskInboxClaimNextAndMarkProcessed(t *testing.T) {
	root := t.TempDir()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	inbox := NewTaskInbox(root, logger)
	if inbox == nil {
		t.Fatalf("expected inbox")
	}
	if err := inbox.EnsureDirs(); err != nil {
		t.Fatalf("EnsureDirs() error = %v", err)
	}

	payload := `{"summary":"alpha task","class":"analysis","primary_modality":"text","prompt":"hello"}`
	if err := os.WriteFile(filepath.Join(root, "001-alpha.json"), []byte(payload), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "zzz-ignore.txt"), []byte("ignore"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	queued, err := inbox.ClaimNext()
	if err != nil {
		t.Fatalf("ClaimNext() error = %v", err)
	}
	if queued == nil {
		t.Fatalf("expected queued task")
	}
	if queued.Task.ID != "001-alpha" {
		t.Fatalf("unexpected task id: %q", queued.Task.ID)
	}
	if queued.Task.PrimaryModality != modality.Text {
		t.Fatalf("unexpected primary modality: %q", queued.Task.PrimaryModality)
	}
	if err := inbox.MarkProcessed(queued); err != nil {
		t.Fatalf("MarkProcessed() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "processed", "001-alpha.json")); err != nil {
		t.Fatalf("expected processed archive file: %v", err)
	}
}

func TestTaskInboxMarkFailedWritesErrorNote(t *testing.T) {
	root := t.TempDir()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	inbox := NewTaskInbox(root, logger)
	if err := inbox.EnsureDirs(); err != nil {
		t.Fatalf("EnsureDirs() error = %v", err)
	}

	if err := os.WriteFile(filepath.Join(root, "002-bad.json"), []byte(`{not-json}`), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	queued, err := inbox.ClaimNext()
	if err == nil {
		t.Fatalf("expected ClaimNext() error")
	}
	if queued == nil {
		t.Fatalf("expected queued task on decode failure")
	}
	if err := inbox.MarkFailed(queued, err); err != nil {
		t.Fatalf("MarkFailed() error = %v", err)
	}

	failedJSON := filepath.Join(root, "failed", "002-bad.json")
	if _, statErr := os.Stat(failedJSON); statErr != nil {
		t.Fatalf("expected failed archive file: %v", statErr)
	}
	errorNote, readErr := os.ReadFile(failedJSON + ".error.txt")
	if readErr != nil {
		t.Fatalf("expected error note: %v", readErr)
	}
	if !strings.Contains(string(errorNote), "decode task file") {
		t.Fatalf("unexpected error note: %s", string(errorNote))
	}
}
