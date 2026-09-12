package cache

import (
	"os"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	err := Set("test-cmd", "Usage: test-cmd [options]")
	if err != nil {
		t.Fatalf("failed to set cache: %v", err)
	}

	doc, found := Get("test-cmd")
	if !found {
		t.Fatal("expected to find cached doc")
	}
	if doc != "Usage: test-cmd [options]" {
		t.Errorf("unexpected cached doc: %s", doc)
	}
}

func TestGetMissing(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	_, found := Get("nonexistent")
	if found {
		t.Error("expected not found for nonexistent command")
	}
}

func TestClear(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	_ = Set("test-cmd", "some doc")
	err := Clear()
	if err != nil {
		t.Fatalf("failed to clear cache: %v", err)
	}

	_, found := Get("test-cmd")
	if found {
		t.Error("expected not found after clear")
	}
}

func TestSize(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	_ = Set("cmd1", "doc1")
	_ = Set("cmd2", "doc2")

	size, err := Size()
	if err != nil {
		t.Fatalf("failed to get size: %v", err)
	}
	if size != 2 {
		t.Errorf("expected size 2, got %d", size)
	}
}
