package analyzer

import (
	"context"
	"fmt"
	"os"

	"github.com/sidx1/sure/internal/ai"
	"github.com/sidx1/sure/internal/config"
	"github.com/sidx1/sure/internal/discovery"
	"github.com/sidx1/sure/internal/parser"
)

// Pipeline defines the pluggable analysis pipeline block.
// Forks and developers can implement custom validation filters,
// heuristic pre-checks, offline rule engines, or multi-model ensembles.
type Pipeline interface {
	Analyze(ctx context.Context, command string) (*ai.AnalysisResult, error)
}

type DefaultPipeline struct{}

func NewDefaultPipeline() *DefaultPipeline {
	return &DefaultPipeline{}
}

func (p *DefaultPipeline) Analyze(ctx context.Context, command string) (*ai.AnalysisResult, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if !cfg.Enabled {
		return &ai.AnalysisResult{Warn: false}, nil
	}

	parsed, err := parser.ParseCommand(command)
	if err != nil {
		return nil, err
	}

	for _, ex := range cfg.ExcludedCommands {
		if parsed.Name == ex {
			return &ai.AnalysisResult{Warn: false}, nil
		}
	}

	for _, ac := range cfg.AlwaysConfirm {
		if command == ac || parsed.Name == ac {
			return &ai.AnalysisResult{
				Warn:       true,
				Confidence: 1.0,
				Category:   "risky",
				Reason:     "Command is in the always-confirm list.",
			}, nil
		}
	}

	provider, err := ai.GetProvider(cfg.Provider)
	if err != nil {
		return nil, err
	}

	doc := discovery.GetCommandDoc(parsed.Name)
	wd, _ := os.Getwd()
	shell := os.Getenv("SHELL")

	req := ai.AnalysisRequest{
		Command:    command,
		CommandDoc: doc,
		WorkingDir: wd,
		Shell:      shell,
	}

	res, err := provider.Analyze(ctx, req)
	if err != nil {
		return nil, err
	}

	if res.Warn {
		if cfg.Sensitivity == "low" && res.Confidence < 0.8 {
			res.Warn = false
		} else if cfg.Sensitivity == "medium" && res.Confidence < 0.5 {
			res.Warn = false
		}
	}

	return res, nil
}

var globalPipeline Pipeline = NewDefaultPipeline()

func SetGlobalPipeline(p Pipeline) {
	if p != nil {
		globalPipeline = p
	}
}

func Analyze(command string) (*ai.AnalysisResult, error) {
	return globalPipeline.Analyze(context.Background(), command)
}
