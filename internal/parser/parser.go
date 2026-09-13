package parser

import (
	"strings"
	"unicode"
)

// Command represents a parsed terminal command.
type Command struct {
	Name   string            `json:"name"`
	Args   []string          `json:"args"`
	Raw    string            `json:"raw"`
	Pipes  []*Command        `json:"pipes,omitempty"`
	Chains []*ChainedCommand `json:"chains,omitempty"`
}

type ChainedCommand struct {
	Operator string   `json:"operator"` // "&&", "||", ";"
	Command  *Command `json:"command"`
}

// Parser defines the pluggable command parser interface block.
// Developers and forks can swap out this block (e.g. tree-sitter-bash, mvdan.cc/sh)
// without breaking the rest of the application.
type Parser interface {
	Parse(input string) (*Command, error)
}

// DefaultParser is the standard tokenizer implementation.
type DefaultParser struct{}

func NewDefaultParser() *DefaultParser {
	return &DefaultParser{}
}

func (p *DefaultParser) Parse(input string) (*Command, error) {
	return ParseCommand(input)
}

var globalParser Parser = NewDefaultParser()

// SetGlobalParser allows forks and extensions to replace the parser block.
func SetGlobalParser(p Parser) {
	if p != nil {
		globalParser = p
	}
}

// ParseCommand parses a shell command line string into a structured Command object.
func ParseCommand(input string) (*Command, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return &Command{Raw: input}, nil
	}

	var tokens []string
	var currentToken strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	var escapeNext bool

	for _, char := range input {
		if escapeNext {
			currentToken.WriteRune(char)
			escapeNext = false
			continue
		}
		if char == '\\' {
			escapeNext = true
			continue
		}
		if char == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}
		if char == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}
		if unicode.IsSpace(char) && !inSingleQuote && !inDoubleQuote {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			continue
		}
		currentToken.WriteRune(char)
	}
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	if len(tokens) == 0 {
		return &Command{Raw: input}, nil
	}

	name := tokens[0]
	args := tokens[1:]

	return &Command{
		Name: name,
		Args: args,
		Raw:  input,
	}, nil
}
