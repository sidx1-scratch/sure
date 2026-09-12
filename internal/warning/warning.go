package warning

import (
	"github.com/sidx1/sure/internal/ai"
	"github.com/sidx1/sure/internal/tui"
)

func Display(result *ai.AnalysisResult, command string) bool {
	if !result.Warn {
		return true
	}
	
	category := result.Category
	if category == "" {
		category = "risky"
	}
	
	return tui.ShowWarning(category, command, result.Reason)
}
