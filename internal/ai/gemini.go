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

type GeminiProvider struct {
	model string
}

func NewGeminiProvider() *GeminiProvider {
	cfg, err := config.Load()
	model := "gemini-2.0-flash"
	if err == nil && cfg.Model != "" {
		model = cfg.Model
	}
	return &GeminiProvider{model: model}
}

func (p *GeminiProvider) Name() string {
	return "gemini"
}

type geminiRequest struct {
	Contents         []geminiContent `json:"contents"`
	SystemInstruction *geminiContent  `json:"systemInstruction,omitempty"`
	GenerationConfig  geminiConfig    `json:"generationConfig"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiConfig struct {
	ResponseMimeType string `json:"responseMimeType"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (p *GeminiProvider) Analyze(ctx context.Context, req AnalysisRequest) (*AnalysisResult, error) {
	apiKey, err := keyring.GetAPIKey("gemini")
	if err != nil {
		return nil, fmt.Errorf("gemini api key not found: %w", err)
	}

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
Return ONLY a JSON object with these exact keys:
{
  "warn": boolean,       // true if there is a real issue to warn about
  "confidence": float,   // 0.0 to 1.0, how confident you are in the warning
  "category": string,    // one of: "destructive", "typo", "suspicious", "invalid_syntax", "risky"
  "reason": string       // concise explanation of WHAT could happen and WHY (1-2 sentences max)
}

The "reason" must explain the actual consequence, not just label it as dangerous.
Example good reason: "'-R' makes the permission change recursive, and '/' targets the entire filesystem."
Example bad reason: "This command is dangerous."

If warn is false, reason can be empty.`
	
	prompt := fmt.Sprintf("Command: %s\nWorking Dir: %s\nShell: %s\nCommand Doc: %s\n", 
		req.Command, req.WorkingDir, req.Shell, req.CommandDoc)

	geminiReq := geminiRequest{
		SystemInstruction: &geminiContent{
			Parts: []geminiPart{{Text: systemPrompt}},
		},
		Contents: []geminiContent{
			{
				Role: "user",
				Parts: []geminiPart{{Text: prompt}},
			},
		},
		GenerationConfig: geminiConfig{
			ResponseMimeType: "application/json",
		},
	}

	body, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", p.model, apiKey)
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gemini API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var geminiResp geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from gemini")
	}

	text := geminiResp.Candidates[0].Content.Parts[0].Text
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var result AnalysisResult
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
