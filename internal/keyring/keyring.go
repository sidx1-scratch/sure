package keyring

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sidx1/sure/internal/config"
	"github.com/zalando/go-keyring"
)

const serviceName = "sure-cli"

func StoreAPIKey(provider, key string) error {
	err := keyring.Set(serviceName, provider, key)
	if err == nil {
		return nil
	}

	// Fallback to file
	dir := config.ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir for fallback keyring: %w", err)
	}
	path := filepath.Join(dir, fmt.Sprintf("%s.key", provider))
	return os.WriteFile(path, []byte(key), 0600)
}

func GetAPIKey(provider string) (string, error) {
	// Check env first
	envMap := map[string]string{
		"gemini": "GEMINI_API_KEY",
		"openai": "OPENAI_API_KEY", // placeholder if they add more
	}
	
	if envKey, ok := envMap[provider]; ok {
		if val := os.Getenv(envKey); val != "" {
			return val, nil
		}
	}
	
	if val := os.Getenv("SURE_API_KEY"); val != "" {
		return val, nil
	}

	// Check keyring
	key, err := keyring.Get(serviceName, provider)
	if err == nil {
		return key, nil
	}

	// Check file
	path := filepath.Join(config.ConfigDir(), fmt.Sprintf("%s.key", provider))
	data, err := os.ReadFile(path)
	if err == nil {
		return string(data), nil
	}

	return "", fmt.Errorf("api key not found for provider %s", provider)
}

func DeleteAPIKey(provider string) error {
	err := keyring.Delete(serviceName, provider)
	
	// Try deleting fallback file too
	path := filepath.Join(config.ConfigDir(), fmt.Sprintf("%s.key", provider))
	_ = os.Remove(path)

	if err != nil {
		return fmt.Errorf("failed to delete from keyring: %w", err)
	}
	return nil
}
