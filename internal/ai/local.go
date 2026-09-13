package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sidx1/sure/internal/config"
	"github.com/sidx1/sure/internal/keyring"
)

// LocalProvider connects to locally running models via Ollama or OpenAI-compatible endpoints
// like LM Studio, LocalAI, vLLM, or text-generation-webui.
// It requires NO internet connection and NO paid API subscription.
type LocalProvider struct {
	endpoint string
	model    string
}

func NewLocalProvider() *LocalProvider {
	cfg, err := config.Load()
	endpoint := "http://localhost:11434"
	model := "llama3.2"

	if err == nil {
		if cfg.Local.Endpoint != "" {
			endpoint = cfg.Local.Endpoint
		}
		if cfg.Local.Model != "" {
			model = cfg.Local.Model
		}
	}
	return &LocalProvider{
		endpoint: strings.TrimRight(endpoint, "/"),
		model:    model,
	}
}

func (p *LocalProvider) Name() string {
	return "local"
}

type openAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatRequest struct {
	Model       string              `json:"model"`
	Messages    []openAIChatMessage `json:"messages"`
	Temperature float64             `json:"temperature"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type ollamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	System string `json:"system"`
	Format string `json:"format"`
	Stream bool   `json:"stream"`
}

type ollamaGenerateResponse struct {
	Response string `json:"response"`
}

func (p *LocalProvider) Analyze(ctx context.Context, req AnalysisRequest) (*AnalysisResult, error) {
	systemPrompt := `You are a terminal command safety analyzer. Your ONLY job is to analyze shell commands and report potential issues.

RULES:
1. NEVER suggest corrections or replacements for commands.
2. NEVER modify the user's command.
3. Only warn when you have a SPECIFIC, CONCRETE reason.
4. Do NOT warn just because a command is unfamiliar or advanced.
5. Do NOT assume destructive = wrong. Advanced users run destructive commands intentionally.
6. If the command looks fine, set "warn" to false.

ANALYZE FOR:
- Syntax errors or invalid subcommands (typos like "git chekcout")
- Potentially destructive operations (rm -rf /, chmod -R 777 /, mkfs on mounted drives)
- Suspicious paths (operating on / or /etc or /boot recursively)
- Contradictory or unusual flag combinations
- Commands that could have unintended side effects
- Nonexistent or suspicious file paths when obvious

RESPONSE FORMAT:
Return ONLY a valid JSON object with these exact keys:
{
  "warn": boolean,
  "confidence": float,
  "category": string,
  "reason": string
}

Category must be one of: "destructive", "typo", "suspicious", "invalid_syntax", "risky".
If warn is false, reason can be empty.`

	userPrompt := fmt.Sprintf("Command: %s\nWorking Dir: %s\nShell: %s\nCommand Doc: %s\n",
		req.Command, req.WorkingDir, req.Shell, req.CommandDoc)

	// Determine if endpoint is Ollama native or OpenAI-compatible
	if strings.Contains(p.endpoint, ":11434") && !strings.Contains(p.endpoint, "/v1") {
		// Try Ollama native endpoint: /api/generate
		res, err := p.callOllamaNative(ctx, systemPrompt, userPrompt)
		if err == nil {
			return res, nil
		}
	}

	// Fallback or explicit OpenAI-compatible endpoint: /v1/chat/completions or /chat/completions
	return p.callOpenAICompatible(ctx, systemPrompt, userPrompt)
}

func (p *LocalProvider) callOllamaNative(ctx context.Context, systemPrompt, userPrompt string) (*AnalysisResult, error) {
	url := fmt.Sprintf("%s/api/generate", p.endpoint)

	reqPayload := ollamaGenerateRequest{
		Model:  p.model,
		System: systemPrompt,
		Prompt: userPrompt,
		Format: "json",
		Stream: false,
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("could not connect to local model at %s (is Ollama running?): %w", p.endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("local model error (%d): %s", resp.StatusCode, string(b))
	}

	var ollamaResp ollamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, err
	}

	return parseJSONResult(ollamaResp.Response)
}

func (p *LocalProvider) callOpenAICompatible(ctx context.Context, systemPrompt, userPrompt string) (*AnalysisResult, error) {
	url := fmt.Sprintf("%s/v1/chat/completions", p.endpoint)
	if strings.HasSuffix(p.endpoint, "/v1") {
		url = fmt.Sprintf("%s/chat/completions", p.endpoint)
	}

	chatReq := openAIChatRequest{
		Model: p.model,
		Messages: []openAIChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.1,
	}

	body, err := json.Marshal(chatReq)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Pass API key if configured
	if key, err := keyring.GetAPIKey("local"); err == nil && key != "" {
		httpReq.Header.Set("Authorization", "Bearer "+key)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("could not connect to local model at %s: %w", p.endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("local model error (%d): %s", resp.StatusCode, string(b))
	}

	var chatResp openAIChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, err
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices from local model")
	}

	return parseJSONResult(chatResp.Choices[0].Message.Content)
}

func parseJSONResult(text string) (*AnalysisResult, error) {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	// Extract JSON substring if surrounded by extra text
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start != -1 && end != -1 && end > start {
		text = text[start : end+1]
	}

	var result AnalysisResult
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI JSON response %q: %w", text, err)
	}

	return &result, nil
}

// OpenAICompatibleProvider handles generic OpenAI-compatible APIs (including local proxies, vLLM, etc.)
type OpenAICompatibleProvider struct {
	providerName string
}

func NewOpenAICompatibleProvider(name string) *OpenAICompatibleProvider {
	return &OpenAICompatibleProvider{providerName: name}
}

func (p *OpenAICompatibleProvider) Name() string {
	return p.providerName
}

func (p *OpenAICompatibleProvider) Analyze(ctx context.Context, req AnalysisRequest) (*AnalysisResult, error) {
	local := NewLocalProvider()
	return local.Analyze(ctx, req)
}
