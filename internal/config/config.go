package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Provider         string   `yaml:"provider"`
	Sensitivity      string   `yaml:"sensitivity"`
	ExcludedCommands []string `yaml:"excluded_commands"`
	AlwaysConfirm    []string `yaml:"always_confirm"`
	AnalyzePipelines bool     `yaml:"analyze_pipelines"`
	CacheDir         string   `yaml:"cache_dir"`
	Enabled          bool     `yaml:"enabled"`
	Model            string   `yaml:"model"`
}

func DefaultConfig() *Config {
	return &Config{
		Provider:         "gemini",
		Sensitivity:      "medium",
		ExcludedCommands: []string{"cd", "ls", "pwd", "echo"},
		AlwaysConfirm:    []string{"rm -rf /", "mkfs"},
		AnalyzePipelines: true,
		CacheDir:         CacheDir(),
		Enabled:          true,
		Model:            "gemini-2.0-flash",
	}
}

func ConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/.sure"
	}
	return filepath.Join(home, ".config", "sure")
}

func CacheDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/.cache/sure"
	}
	return filepath.Join(home, ".cache", "sure")
}

func Load() (*Config, error) {
	path := filepath.Join(ConfigDir(), "config.yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

func Save(cfg *Config) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
