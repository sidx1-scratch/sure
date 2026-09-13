package stats

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStatsTracker(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	tracker := &FileTracker{
		filePath: filepath.Join(tmpDir, "events.jsonl"),
	}

	// 1. Initial summary should be clean
	s, err := tracker.GetSummary()
	if err != nil {
		t.Fatalf("failed to get summary: %v", err)
	}
	if s.TotalAnalyzed != 0 || s.AccuracyScore != 100.0 {
		t.Errorf("expected 0 analyzed and 100%% accuracy, got %d and %f", s.TotalAnalyzed, s.AccuracyScore)
	}

	// 2. Record 3 safe commands
	for i := 0; i < 3; i++ {
		_ = tracker.Record(CommandEvent{
			Timestamp: time.Now(),
			Command:   "git status",
			Warned:    false,
			Proceeded: true,
		})
	}

	// 3. Record 1 caught mistake (user warned, cancelled)
	_ = tracker.Record(CommandEvent{
		Timestamp:  time.Now(),
		Command:    "rm -rf /",
		Warned:     true,
		Proceeded:  false,
		Category:   "destructive",
		Reason:     "Deletes entire filesystem",
		Confidence: 0.99,
	})

	s, err = tracker.GetSummary()
	if err != nil {
		t.Fatalf("failed to get summary: %v", err)
	}

	if s.TotalAnalyzed != 4 {
		t.Errorf("expected 4 analyzed, got %d", s.TotalAnalyzed)
	}
	if s.MistakesCaught != 1 {
		t.Errorf("expected 1 mistake caught, got %d", s.MistakesCaught)
	}
	if s.AccuracyScore != 75.0 {
		t.Errorf("expected 75%% accuracy, got %f", s.AccuracyScore)
	}
	if len(s.RecentMistakes) != 1 || s.RecentMistakes[0].Command != "rm -rf /" {
		t.Errorf("expected recent mistake 'rm -rf /', got %+v", s.RecentMistakes)
	}
}
