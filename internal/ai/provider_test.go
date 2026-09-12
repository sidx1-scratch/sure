package ai

import (
	"encoding/json"
	"testing"
)

func TestAnalysisResultMarshal(t *testing.T) {
	result := AnalysisResult{
		Warn:       true,
		Confidence: 0.91,
		Category:   "destructive",
		Reason:     "The command recursively changes permissions on the entire filesystem.",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled AnalysisResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if unmarshaled.Warn != true {
		t.Error("expected warn to be true")
	}
	if unmarshaled.Confidence != 0.91 {
		t.Errorf("expected confidence 0.91, got %f", unmarshaled.Confidence)
	}
	if unmarshaled.Category != "destructive" {
		t.Errorf("expected category 'destructive', got '%s'", unmarshaled.Category)
	}
}

func TestGetProviderGemini(t *testing.T) {
	provider, err := GetProvider("gemini")
	if err != nil {
		t.Fatalf("failed to get gemini provider: %v", err)
	}
	if provider.Name() != "gemini" {
		t.Errorf("expected provider name 'gemini', got '%s'", provider.Name())
	}
}

func TestGetProviderUnknown(t *testing.T) {
	_, err := GetProvider("nonexistent")
	if err == nil {
		t.Error("expected error for unknown provider")
	}
}
