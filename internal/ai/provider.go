package ai

import (
	"context"
	"fmt"
	"sync"
)

type AnalysisRequest struct {
	Command    string `json:"command"`
	CommandDoc string `json:"command_doc"`
	WorkingDir string `json:"working_dir"`
	Shell      string `json:"shell"`
}

type AnalysisResult struct {
	Warn       bool    `json:"warn"`
	Confidence float64 `json:"confidence"`
	Category   string  `json:"category"` // "destructive", "typo", "suspicious", "invalid_syntax", "risky"
	Reason     string  `json:"reason"`
}

// Provider is the modular AI block interface.
// Developers and forks can easily replace or add custom providers
// by implementing this interface and calling RegisterProvider().
type Provider interface {
	Name() string
	Analyze(ctx context.Context, req AnalysisRequest) (*AnalysisResult, error)
}

var (
	registryMu sync.RWMutex
	providers  = make(map[string]ProviderFactory)
)

// ProviderFactory allows lazy creation of providers based on name
type ProviderFactory func() (Provider, error)

// RegisterProvider registers an AI provider block in the registry
func RegisterProvider(name string, factory ProviderFactory) {
	registryMu.Lock()
	defer registryMu.Unlock()
	providers[name] = factory
}

// GetProvider retrieves a provider from the registry
func GetProvider(name string) (Provider, error) {
	registryMu.RLock()
	factory, ok := providers[name]
	registryMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown provider %q (registered: %v)", name, AvailableProviders())
	}
	return factory()
}

// AvailableProviders returns a list of all registered provider names
func AvailableProviders() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()
	names := make([]string, 0, len(providers))
	for name := range providers {
		names = append(names, name)
	}
	return names
}

func init() {
	// Register default modular blocks
	RegisterProvider("gemini", func() (Provider, error) {
		return NewGeminiProvider(), nil
	})
	RegisterProvider("local", func() (Provider, error) {
		return NewLocalProvider(), nil
	})
	RegisterProvider("ollama", func() (Provider, error) {
		return NewLocalProvider(), nil
	})
	RegisterProvider("openai", func() (Provider, error) {
		return NewOpenAICompatibleProvider("openai"), nil
	})
}
