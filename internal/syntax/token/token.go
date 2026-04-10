// Package token defines the lexical tokens for shell script parsing.
package token

import (
	"fmt"
)

// TokenType represents the type of a token.
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF
	COMMENT
	NEWLINE
	WHITESPACE

	// Literals
	LITERAL
	SINGLE_QUOTED
	DOUBLE_QUOTED

	// Operators and punctuation
	ASSIGN        // =
	PIPE          // |
	SEMICOLON     // ;
	AND           // &&
	OR            // ||
	BACKGROUND    // &
	REDIRECT_IN   // <
	REDIRECT_OUT  // >
	REDIRECT_APPEND // >>
	REDIRECT_HERE // <<

	// Variable expansions
	DOLLAR        // $
	DOLLAR_BRACE  // ${
	DOLLAR_PAREN  // $(
	DOLLAR_DOUBLE_PAREN // $((
	DOLLAR_BRACKET // $[

	// Braces and brackets
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]

	// Keywords
	IF
	THEN
	ELSE
	ELIF
	FI
	FOR
	IN
	DO
	DONE
	WHILE
	UNTIL
	CASE
	ESAC
	SELECT
	FUNCTION
)

var tokenNames = map[TokenType]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",
	NEWLINE: "NEWLINE",
	WHITESPACE: "WHITESPACE",

	LITERAL:       "LITERAL",
	SINGLE_QUOTED: "SINGLE_QUOTED",
	DOUBLE_QUOTED: "DOUBLE_QUOTED",

	ASSIGN:        "=",
	PIPE:          "|",
	SEMICOLON:     ";",
	AND:           "&&",
	OR:            "||",
	BACKGROUND:    "&",
	REDIRECT_IN:   "<",
	REDIRECT_OUT:  ">",
	REDIRECT_APPEND: ">>",
	REDIRECT_HERE: "<<",

	DOLLAR:        "$",
	DOLLAR_BRACE:  "${",
	DOLLAR_PAREN:  "$(",
	DOLLAR_DOUBLE_PAREN: "$((",
	DOLLAR_BRACKET: "$[",

	LPAREN:    "(",
	RPAREN:    ")",
	LBRACE:    "{",
	RBRACE:    "}",
	LBRACKET:  "[",
	RBRACKET:  "]",

	IF:       "if",
	THEN:     "then",
	ELSE:     "else",
	ELIF:     "elif",
	FI:       "fi",
	FOR:      "for",
	IN:       "in",
	DO:       "do",
	DONE:     "done",
	WHILE:    "while",
	UNTIL:    "until",
	CASE:     "case",
	ESAC:     "esac",
	SELECT:   "select",
	FUNCTION: "function",
}

// String returns the string representation of the token type.
func (t TokenType) String() string {
	if name, ok := tokenNames[t]; ok {
		return name
	}
	return fmt.Sprintf("TokenType(%d)", t)
}

// Position represents a location in the source code.
type Position struct {
	Filename string // filename, if any
	Offset   int    // byte offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (byte count)
}

// IsValid returns true if the position is valid.
func (pos *Position) IsValid() bool {
	return pos.Line > 0
}

// String returns a string representation of the position.
func (pos Position) String() string {
	s := pos.Filename
	if pos.IsValid() {
		if s != "" {
			s += ":"
		}
		s += fmt.Sprintf("%d:%d", pos.Line, pos.Column)
	}
	if s == "" {
		s = "-"
	}
	return s
}

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Literal string
	Pos     Position
}

// New creates a new token.
func New(typ TokenType, lit string, pos Position) Token {
	return Token{
		Type:    typ,
		Literal: lit,
		Pos:     pos,
	}
}

// String returns a string representation of the token.
func (t Token) String() string {
	return fmt.Sprintf("%s %q at %s", t.Type, t.Literal, t.Pos)
}
