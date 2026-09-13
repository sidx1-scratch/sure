package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sidx1/sure/internal/analyzer"
	"github.com/sidx1/sure/internal/stats"
	"github.com/sidx1/sure/internal/warning"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze a command",
	Run: func(cmd *cobra.Command, args []string) {
		command, _ := cmd.Flags().GetString("command")
		jsonOut, _ := cmd.Flags().GetBool("json")

		if command == "" {
			if jsonOut {
				fmt.Println(`{"error": "no command provided"}`)
			} else {
				fmt.Println("Error: no command provided")
			}
			os.Exit(1)
		}

		result, err := analyzer.Analyze(command)
		if err != nil {
			if jsonOut {
				fmt.Printf(`{"error": "%s"}\n`, err.Error())
			} else {
				fmt.Fprintf(os.Stderr, "Analysis error: %v\n", err)
			}
			// Safe default: let the command run
			_ = stats.Record(stats.CommandEvent{
				Command:   command,
				Warned:    false,
				Proceeded: true,
			})
			os.Exit(0)
		}

		if jsonOut {
			data, _ := json.Marshal(result)
			fmt.Println(string(data))
			if result.Warn {
				_ = stats.Record(stats.CommandEvent{
					Command:    command,
					Warned:     true,
					Proceeded:  false,
					Category:   result.Category,
					Reason:     result.Reason,
					Confidence: result.Confidence,
				})
				os.Exit(1)
			}
			_ = stats.Record(stats.CommandEvent{
				Command:   command,
				Warned:    false,
				Proceeded: true,
			})
			os.Exit(0)
		}

		if !result.Warn {
			// Clean execution! Record safe command for typing accuracy tracking
			_ = stats.Record(stats.CommandEvent{
				Command:   command,
				Warned:    false,
				Proceeded: true,
			})
			os.Exit(0)
		}

		// Show interactive warning
		proceed := warning.Display(result, command)

		// Record outcome into the stats tracker block!
		_ = stats.Record(stats.CommandEvent{
			Command:    command,
			Warned:     true,
			Proceeded:  proceed,
			Category:   result.Category,
			Reason:     result.Reason,
			Confidence: result.Confidence,
		})

		if !proceed {
			// User chose NO -> Mistake caught!
			os.Exit(1)
		}
		// User chose YES -> Proceed with command
		os.Exit(0)
	},
}

func init() {
	analyzeCmd.Flags().String("command", "", "Command to analyze")
	analyzeCmd.Flags().Bool("json", false, "Output JSON result")
	rootCmd.AddCommand(analyzeCmd)
}
