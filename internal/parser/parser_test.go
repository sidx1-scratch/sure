package parser

import (
	"testing"
)

func TestParseCommand_Simple(t *testing.T) {
	cmd, err := ParseCommand("ls -la")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Name != "ls" {
		t.Errorf("expected name 'ls', got '%s'", cmd.Name)
	}
	if len(cmd.Args) != 1 || cmd.Args[0] != "-la" {
		t.Errorf("expected args ['-la'], got %v", cmd.Args)
	}
}

func TestParseCommand_QuotedArgs(t *testing.T) {
	cmd, err := ParseCommand(`echo "hello world"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Name != "echo" {
		t.Errorf("expected name 'echo', got '%s'", cmd.Name)
	}
	if len(cmd.Args) != 1 || cmd.Args[0] != "hello world" {
		t.Errorf("expected args ['hello world'], got %v", cmd.Args)
	}
}

func TestParseCommand_SingleQuotes(t *testing.T) {
	cmd, err := ParseCommand("echo 'hello world'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cmd.Args) != 1 || cmd.Args[0] != "hello world" {
		t.Errorf("expected args ['hello world'], got %v", cmd.Args)
	}
}

func TestParseCommand_EscapedChars(t *testing.T) {
	cmd, err := ParseCommand(`echo hello\ world`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cmd.Args) != 1 || cmd.Args[0] != "hello world" {
		t.Errorf("expected args ['hello world'], got %v", cmd.Args)
	}
}

func TestParseCommand_Empty(t *testing.T) {
	cmd, err := ParseCommand("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Name != "" {
		t.Errorf("expected empty name, got '%s'", cmd.Name)
	}
}

func TestParseCommand_MultipleArgs(t *testing.T) {
	cmd, err := ParseCommand("git commit -m 'initial commit'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Name != "git" {
		t.Errorf("expected name 'git', got '%s'", cmd.Name)
	}
	if len(cmd.Args) != 3 {
		t.Fatalf("expected 3 args, got %d: %v", len(cmd.Args), cmd.Args)
	}
	if cmd.Args[0] != "commit" {
		t.Errorf("expected first arg 'commit', got '%s'", cmd.Args[0])
	}
	if cmd.Args[1] != "-m" {
		t.Errorf("expected second arg '-m', got '%s'", cmd.Args[1])
	}
	if cmd.Args[2] != "initial commit" {
		t.Errorf("expected third arg 'initial commit', got '%s'", cmd.Args[2])
	}
}

func TestParseCommand_RawPreserved(t *testing.T) {
	raw := "rm -rf ~/Downloads/*"
	cmd, err := ParseCommand(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Raw != raw {
		t.Errorf("expected raw '%s', got '%s'", raw, cmd.Raw)
	}
}

func TestParseCommand_ComplexFlags(t *testing.T) {
	cmd, err := ParseCommand("chmod -R 777 /")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Name != "chmod" {
		t.Errorf("expected name 'chmod', got '%s'", cmd.Name)
	}
	if len(cmd.Args) != 3 {
		t.Fatalf("expected 3 args, got %d: %v", len(cmd.Args), cmd.Args)
	}
}
