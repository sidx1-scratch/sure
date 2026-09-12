package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Provider != "gemini" {
		t.Errorf("expected default provider 'gemini', got '%s'", cfg.Provider)
	}
	if cfg.Sensitivity != "medium" {
		t.Errorf("expected default sensitivity 'medium', got '%s'", cfg.Sensitivity)
	}
	if !cfg.Enabled {
		t.Error("expected enabled to be true by default")
	}
	if !cfg.AnalyzePipelines {
		t.Error("expected analyze_pipelines to be true by default")
	}
	if cfg.Model != "gemini-2.0-flash" {
		t.Errorf("expected default model 'gemini-2.0-flash', got '%s'", cfg.Model)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg := DefaultConfig()
	cfg.Provider = "test-provider"
	cfg.Sensitivity = "high"

	err := Save(cfg)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Verify file exists
	path := filepath.Join(tmpDir, ".config", "sure", "config.yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("config file was not created")
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if loaded.Provider != "test-provider" {
		t.Errorf("expected provider 'test-provider', got '%s'", loaded.Provider)
	}
	if loaded.Sensitivity != "high" {
		t.Errorf("expected sensitivity 'high', got '%s'", loaded.Sensitivity)
	}
}

func TestLoadMissing(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load should not error for missing config: %v", err)
	}
	if cfg.Provider != "gemini" {
		t.Errorf("expected default provider when config missing, got '%s'", cfg.Provider)
	}
}
