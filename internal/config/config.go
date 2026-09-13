package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type LocalModelConfig struct {
	Endpoint string `yaml:"endpoint"` // e.g. http://localhost:11434/v1 for Ollama, http://localhost:1234/v1 for LM Studio
	Model    string `yaml:"model"`    // e.g. llama3.2, mistral, qwen2.5-coder, etc.
}

type Config struct {
	Provider         string           `yaml:"provider"`           // "gemini", "ollama", "local", "openai"
	Sensitivity      string           `yaml:"sensitivity"`        // "low", "medium", "high"
	ExcludedCommands []string         `yaml:"excluded_commands"`  // commands to skip
	AlwaysConfirm    []string         `yaml:"always_confirm"`     // always warn
	AnalyzePipelines bool             `yaml:"analyze_pipelines"`  // default true
	CacheDir         string           `yaml:"cache_dir"`          // default ~/.cache/sure
	Enabled          bool             `yaml:"enabled"`            // default true
	Model            string           `yaml:"model"`              // default "gemini-2.5-flash" (free-tier eligible)
	Local            LocalModelConfig `yaml:"local"`              // Local model settings (Ollama / OpenAI-compatible)
	WebPort          int              `yaml:"web_port"`           // Web dashboard port (default 7873)
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
		Model:            "gemini-2.5-flash", // Free tier eligible on Google AI Studio
		Local: LocalModelConfig{
			Endpoint: "http://localhost:11434",
			Model:    "llama3.2",
		},
		WebPort: 7873,
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

func DataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/.local/share/sure"
	}
	return filepath.Join(home, ".local", "share", "sure")
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
