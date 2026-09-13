package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sidx1/sure/internal/ai"
	"github.com/sidx1/sure/internal/cache"
	"github.com/sidx1/sure/internal/config"
	"github.com/sidx1/sure/internal/keyring"
	"github.com/sidx1/sure/internal/stats"
	"github.com/sidx1/sure/internal/tui"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose issues with AI providers, local models, keyring, and shell integration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n🩺 sure doctor — System & Environment Diagnostics")
		fmt.Println("=================================================")

		cfg, err := config.Load()
		if err != nil {
			tui.PrintError("Config: Failed to load")
		} else {
			tui.PrintSuccess(fmt.Sprintf("Config: Loaded ok (Provider: %s, Model: %s)", cfg.Provider, cfg.Model))
		}

		// Check AI Provider
		switch cfg.Provider {
		case "gemini":
			key, err := keyring.GetAPIKey("gemini")
			if err != nil || key == "" {
				tui.PrintError("AI Provider (Gemini): Missing API key (run 'sure setup')")
			} else {
				tui.PrintSuccess("AI Provider (Gemini): API key found (Free-tier eligible model)")
			}
		case "local", "ollama":
			tui.PrintInfo(fmt.Sprintf("AI Provider (Local): Testing connection to %s...", cfg.Local.Endpoint))
			client := &http.Client{Timeout: 3 * time.Second}
			resp, err := client.Get(cfg.Local.Endpoint)
			if err != nil {
				tui.PrintError(fmt.Sprintf("AI Provider (Local): Could not reach %s. Is Ollama/LM Studio running?", cfg.Local.Endpoint))
			} else {
				resp.Body.Close()
				tui.PrintSuccess(fmt.Sprintf("AI Provider (Local): Successfully reached local model endpoint (%s, model: %s)", cfg.Local.Endpoint, cfg.Local.Model))
			}
		default:
			tui.PrintInfo(fmt.Sprintf("AI Provider: %s (registered providers: %v)", cfg.Provider, ai.AvailableProviders()))
		}

		// Cache check
		size, err := cache.Size()
		if err != nil {
			tui.PrintError("Doc Cache: Unreachable")
		} else {
			tui.PrintSuccess(fmt.Sprintf("Doc Cache: %d cached command docs at %s", size, config.CacheDir()))
		}

		// Stats Tracker check
		summary, err := stats.GetSummary()
		if err != nil {
			tui.PrintError("Stats Tracker: Failed to read stats")
		} else {
			tui.PrintSuccess(fmt.Sprintf("Stats Tracker: Active (%d mistakes caught, %.1f%% accuracy)", summary.MistakesCaught, summary.AccuracyScore))
		}

		// Shell Integration check
		if os.Getenv("SURE_SKIP") != "" {
			tui.PrintSuccess("Shell Hook: Active in current terminal session")
		} else {
			tui.PrintInfo("Shell Hook: Not detected in current subshell (add eval \"$(sure shell-init --shell bash)\" to .bashrc)")
		}

		fmt.Println("")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
