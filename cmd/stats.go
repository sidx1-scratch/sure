package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sidx1/sure/internal/stats"
	"github.com/sidx1/sure/internal/tui"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show command accuracy, mistakes caught, and typing improvement",
	Run: func(cmd *cobra.Command, args []string) {
		jsonOut, _ := cmd.Flags().GetBool("json")
		badgeOnly, _ := cmd.Flags().GetBool("badge")
		reset, _ := cmd.Flags().GetBool("reset")

		if reset {
			if err := stats.DefaultTracker().Reset(); err != nil {
				tui.PrintError(fmt.Sprintf("Failed to reset statistics: %v", err))
				os.Exit(1)
			}
			tui.PrintSuccess("Command safety statistics reset.")
			return
		}

		summary, err := stats.GetSummary()
		if err != nil {
			tui.PrintError(fmt.Sprintf("Failed to load stats: %v", err))
			os.Exit(1)
		}

		if badgeOnly {
			fmt.Printf("🛡️ %d mistakes caught | %.1f%% accuracy\n", summary.MistakesCaught, summary.AccuracyScore)
			return
		}

		if jsonOut {
			data, _ := json.MarshalIndent(summary, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Println("\n📊 sure — Terminal Command Safety & Typing Improvement")
		fmt.Println("=======================================================")
		fmt.Printf("Total Commands Analyzed:    %d\n", summary.TotalAnalyzed)
		fmt.Printf("🛑 Mistakes Caught:         %d (intercepted & avoided)\n", summary.MistakesCaught)
		fmt.Printf("⚠️  Warnings Overridden:     %d (user confirmed Y)\n", summary.WarningsProceeded)
		fmt.Printf("✨ Safe Commands:           %d (clean execution)\n", summary.SafeAllowed)
		fmt.Printf("🎯 Typing Accuracy Score:   %.1f%%\n", summary.AccuracyScore)
		fmt.Printf("📈 Improvement Trend:       %s\n", summary.ImprovementTrend)

		if len(summary.CategoryBreakdown) > 0 {
			fmt.Println("\nBreakdown by Warning Category:")
			for cat, count := range summary.CategoryBreakdown {
				fmt.Printf("  • %-16s %d\n", cat, count)
			}
		}

		if len(summary.RecentMistakes) > 0 {
			fmt.Println("\nRecent Intercepted Mistakes:")
			for i, ev := range summary.RecentMistakes {
				fmt.Printf("  %d. [%s] %s (%s)\n", i+1, ev.Category, ev.Command, ev.Reason)
			}
		}

		fmt.Println("\n💡 Tip: Run 'sure web' to launch the live web dashboard or embed your safety badge!")
	},
}

func init() {
	statsCmd.Flags().Bool("json", false, "Output JSON summary")
	statsCmd.Flags().Bool("badge", false, "Print compact badge string")
	statsCmd.Flags().Bool("reset", false, "Reset statistics history")
	rootCmd.AddCommand(statsCmd)
}
