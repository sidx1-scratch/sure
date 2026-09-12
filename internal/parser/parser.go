package parser

import (
	"strings"
	"unicode"
)

type Command struct {
	Name   string
	Args   []string
	Raw    string
	Pipes  []*Command
	Chains []*ChainedCommand
}

type ChainedCommand struct {
	Operator string
	Command  *Command
}

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
