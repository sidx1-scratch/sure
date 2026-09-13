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
	ColorBold   = "\033[1m"
)

// UI defines the pluggable user interaction block.
type UI interface {
	SetupWizard() (string, error)
	ShowWarning(category, command, reason string) bool
	PrintSuccess(msg string)
	PrintError(msg string)
	PrintInfo(msg string)
}

type TerminalUI struct{}

func NewTerminalUI() *TerminalUI {
	return &TerminalUI{}
}

func (t *TerminalUI) PrintSuccess(msg string) {
	fmt.Printf("%s✔ %s%s\n", ColorGreen, msg, ColorReset)
}

func (t *TerminalUI) PrintError(msg string) {
	fmt.Printf("%s✖ %s%s\n", ColorRed, msg, ColorReset)
}

func (t *TerminalUI) PrintInfo(msg string) {
	fmt.Printf("%sℹ %s%s\n", ColorCyan, msg, ColorReset)
}

func (t *TerminalUI) SetupWizard() (string, error) {
	return SetupWizard()
}

func (t *TerminalUI) ShowWarning(category, command, reason string) bool {
	return ShowWarning(category, command, reason)
}

var globalUI UI = NewTerminalUI()

func SetGlobalUI(u UI) {
	if u != nil {
		globalUI = u
	}
}

func PrintSuccess(msg string) {
	globalUI.PrintSuccess(msg)
}

func PrintError(msg string) {
	globalUI.PrintError(msg)
}

func PrintInfo(msg string) {
	globalUI.PrintInfo(msg)
}

func SetupWizard() (string, error) {
	fmt.Println("╔═══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                             SURE                                  ║")
	fmt.Println("║                                                                   ║")
	fmt.Println("║               Your terminal's second opinion.                     ║")
	fmt.Println("║                                                                   ║")
	fmt.Println("║  Gemini API key (Free-tier eligible, press Enter for local/skip): ║")
	fmt.Print("║  > ")

	reader := bufio.NewReader(os.Stdin)
	key, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	key = strings.TrimSpace(key)

	fmt.Println("║                                                                   ║")
	fmt.Println("║                         [ Continue ]                              ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════════╝")

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
