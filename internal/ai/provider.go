package ai

import (
	"context"
	"fmt"
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
	Category   string  `json:"category"`
	Reason     string  `json:"reason"`
}

type Provider interface {
	Name() string
	Analyze(ctx context.Context, req AnalysisRequest) (*AnalysisResult, error)
}

func GetProvider(name string) (Provider, error) {
	switch name {
	case "gemini":
		return NewGeminiProvider(), nil
	default:
		return nil, fmt.Errorf("unknown provider %s", name)
	}
}
