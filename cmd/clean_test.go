package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResetDirRecreatesDirectory(t *testing.T) {
	target := filepath.Join(t.TempDir(), "output")

	if err := os.MkdirAll(target, 0755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}

	filePath := filepath.Join(target, "stale.txt")
	if err := os.WriteFile(filePath, []byte("old"), 0644); err != nil {
		t.Fatalf("write stale file: %v", err)
	}

	if err := resetDir(target); err != nil {
		t.Fatalf("reset dir: %v", err)
	}

	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("stat target: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("target is not a directory")
	}

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Fatalf("stale file still exists, err=%v", err)
	}
}

func TestEnsureCleanPathRejectsUnexpectedDirectory(t *testing.T) {
	target := filepath.Join(t.TempDir(), "unexpected")

	err := ensureCleanPath(target)
	if err == nil {
		t.Fatalf("expected error for unexpected directory")
	}
}
