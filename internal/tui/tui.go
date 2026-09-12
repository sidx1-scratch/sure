package tui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
)

func PrintSuccess(msg string) {
	fmt.Printf("%s✔ %s%s\n", ColorGreen, msg, ColorReset)
}

func PrintError(msg string) {
	fmt.Printf("%s✖ %s%s\n", ColorRed, msg, ColorReset)
}

func PrintInfo(msg string) {
	fmt.Printf("%sℹ %s%s\n", ColorCyan, msg, ColorReset)
}

func SetupWizard() (string, error) {
	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║              SURE                    ║")
	fmt.Println("║                                      ║")
	fmt.Println("║    Your terminal's second opinion.   ║")
	fmt.Println("║                                      ║")
	fmt.Println("║    Gemini API key:                   ║")
	fmt.Print("║    > ")
	
	reader := bufio.NewReader(os.Stdin)
	key, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	key = strings.TrimSpace(key)
	
	fmt.Println("║                                      ║")
	fmt.Println("║              [ Continue ]            ║")
	fmt.Println("╚══════════════════════════════════════╝")
	
	return key, nil
}

func ShowWarning(category, command, reason string) bool {
	fmt.Printf("\n%s⚠ This command may be %s.%s\n\n", ColorYellow, category, ColorReset)
	fmt.Printf("Command:\n  %s\n\n", command)
	fmt.Printf("Why:\n  %s\n\n", reason)
	fmt.Printf("Do you want to proceed?\n\n")
	fmt.Printf("[Y] Yes   [N] No\n")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return false
		}
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" || input == "yes" {
			return true
		}
		if input == "n" || input == "no" {
			return false
		}
	}
}
