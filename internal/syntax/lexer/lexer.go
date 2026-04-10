// Package lexer implements a lexical analyzer for shell scripts.
package lexer

import (
	"unicode/utf8"

	"github.com/afeldman/goshellcheck/internal/syntax/token"
)

// Lexer represents a lexical scanner.
type Lexer struct {
	input        string     // input string
	position     int        // current position in input (points to current rune)
	readPosition int        // current reading position in input (after current rune)
	ch           rune       // current rune under examination
	filename     string     // source filename
	line         int        // current line number
	column       int        // current column number
}

// New creates a new lexer for the given input.
func New(input, filename string) *Lexer {
	l := &Lexer{
		input:    input,
		filename: filename,
		line:     1,
		column:   1,
	}
	l.readChar()
	return l
}

// readChar reads the next character from the input.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for "NUL" character
	} else {
		l.ch, _ = utf8.DecodeRuneInString(l.input[l.readPosition:])
	}
	l.position = l.readPosition
	l.readPosition += utf8.RuneLen(l.ch)
	l.column++
}

// peekChar returns the next character without advancing.
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	ch, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return ch
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// Skip whitespace (but not newlines)
	l.skipWhitespace()

	pos := token.Position{
		Filename: l.filename,
		Offset:   l.position,
		Line:     l.line,
		Column:   l.column - 1, // column was incremented in readChar
	}

	switch l.ch {
	case 0:
		tok = token.New(token.EOF, "", pos)
	case '\n':
		tok = token.New(token.NEWLINE, "\n", pos)
		l.readChar()
		l.line++
		l.column = 1
		return tok
	case '#':
		tok = l.readComment()
		return tok
	case ';':
		tok = token.New(token.SEMICOLON, ";", pos)
		l.readChar()
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.New(token.OR, literal, pos)
			l.readChar()
		} else {
			tok = token.New(token.PIPE, "|", pos)
			l.readChar()
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.New(token.AND, literal, pos)
			l.readChar()
		} else {
			tok = token.New(token.BACKGROUND, "&", pos)
			l.readChar()
		}
	case '<':
		if l.peekChar() == '<' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.New(token.REDIRECT_HERE, literal, pos)
			l.readChar()
		} else {
			tok = token.New(token.REDIRECT_IN, "<", pos)
			l.readChar()
		}
	case '>':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.New(token.REDIRECT_APPEND, literal, pos)
			l.readChar()
		} else {
			tok = token.New(token.REDIRECT_OUT, ">", pos)
			l.readChar()
		}
	case '=':
		tok = token.New(token.ASSIGN, "=", pos)
		l.readChar()
	case '(':
		tok = token.New(token.LPAREN, "(", pos)
		l.readChar()
	case ')':
		tok = token.New(token.RPAREN, ")", pos)
		l.readChar()
	case '{':
		tok = token.New(token.LBRACE, "{", pos)
		l.readChar()
	case '}':
		tok = token.New(token.RBRACE, "}", pos)
		l.readChar()
	case '[':
		tok = token.New(token.LBRACKET, "[", pos)
		l.readChar()
	case ']':
		tok = token.New(token.RBRACKET, "]", pos)
		l.readChar()
	case '$':
		tok = l.readDollar()
	default:
		if isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' || l.ch == '-' || l.ch == '.' || l.ch == '/' {
			tok = l.readWord()
		} else if l.ch == '\'' {
			tok = l.readSingleQuoted()
		} else if l.ch == '"' {
			tok = l.readDoubleQuoted()
		} else {
			tok = token.New(token.ILLEGAL, string(l.ch), pos)
			l.readChar()
		}
	}

	return tok
}

// skipWhitespace skips whitespace characters (space, tab, carriage return).
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

// readComment reads a comment token.
func (l *Lexer) readComment() token.Token {
	pos := token.Position{
		Filename: l.filename,
		Offset:   l.position,
		Line:     l.line,
		Column:   l.column - 1,
	}
	
	start := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	
	literal := l.input[start:l.position]
	return token.New(token.COMMENT, literal, pos)
}

// readWord reads a word token (identifier or literal).
func (l *Lexer) readWord() token.Token {
	pos := token.Position{
		Filename: l.filename,
		Offset:   l.position,
		Line:     l.line,
		Column:   l.column - 1,
	}
	
	start := l.position
	for isWordChar(l.ch) {
		l.readChar()
	}
	
	literal := l.input[start:l.position]
	
	// Check if it's a keyword
	if typ, ok := keywords[literal]; ok {
		return token.New(typ, literal, pos)
	}
	
	return token.New(token.LITERAL, literal, pos)
}

// readSingleQuoted reads a single-quoted string.
func (l *Lexer) readSingleQuoted() token.Token {
	pos := token.Position{
		Filename: l.filename,
		Offset:   l.position,
		Line:     l.line,
		Column:   l.column - 1,
	}
	
	start := l.position
	l.readChar() // skip opening quote
	
	for l.ch != '\'' && l.ch != 0 {
		if l.ch == '\n' {
			l.line++
			l.column = 1
		}
		l.readChar()
	}
	
	if l.ch == '\'' {
		l.readChar() // skip closing quote
	}
	
	literal := l.input[start:l.position]
	return token.New(token.SINGLE_QUOTED, literal, pos)
}

// readDoubleQuoted reads a double-quoted string.
func (l *Lexer) readDoubleQuoted() token.Token {
	pos := token.Position{
		Filename: l.filename,
		Offset:   l.position,
		Line:     l.line,
		Column:   l.column - 1,
	}
	
	start := l.position
	l.readChar() // skip opening quote
	
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar() // skip backslash
		}
		if l.ch == '\n' {
			l.line++
			l.column = 1
		}
		l.readChar()
	}
	
	if l.ch == '"' {
		l.readChar() // skip closing quote
	}
	
	literal := l.input[start:l.position]
	return token.New(token.DOUBLE_QUOTED, literal, pos)
}

// readDollar reads a dollar token and its possible expansions.
func (l *Lexer) readDollar() token.Token {
	pos := token.Position{
		Filename: l.filename,
		Offset:   l.position,
		Line:     l.line,
		Column:   l.column - 1,
	}
	
	l.readChar() // skip '$'
	
	switch l.ch {
	case '{':
		l.readChar()
		return token.New(token.DOLLAR_BRACE, "${", pos)
	case '(':
		l.readChar()
		if l.ch == '(' {
			l.readChar()
			return token.New(token.DOLLAR_DOUBLE_PAREN, "$((", pos)
		}
		return token.New(token.DOLLAR_PAREN, "$(", pos)
	case '[':
		l.readChar()
		return token.New(token.DOLLAR_BRACKET, "$[", pos)
	default:
		// Simple variable expansion like $var
		return token.New(token.DOLLAR, "$", pos)
	}
}

// isLetter checks if a rune is a letter.
func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit checks if a rune is a digit.
func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// isWordChar checks if a rune can be part of a word.
func isWordChar(ch rune) bool {
	return isLetter(ch) || isDigit(ch) || ch == '_' || ch == '-' || ch == '.' || ch == '/' ||
		ch == ':' || ch == '+' || ch == '*' || ch == '?' || ch == '!' || ch == '@' ||
		ch == '#' || ch == '%' || ch == '^' || ch == '~' || ch == '='
}

// keywords maps keyword strings to their token types.
var keywords = map[string]token.TokenType{
	"if":       token.IF,
	"then":     token.THEN,
	"else":     token.ELSE,
	"elif":     token.ELIF,
	"fi":       token.FI,
	"for":      token.FOR,
	"in":       token.IN,
	"do":       token.DO,
	"done":     token.DONE,
	"while":    token.WHILE,
	"until":    token.UNTIL,
	"case":     token.CASE,
	"esac":     token.ESAC,
	"select":   token.SELECT,
	"function": token.FUNCTION,
}
