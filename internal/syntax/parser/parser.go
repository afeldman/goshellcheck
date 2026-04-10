// Package parser implements a parser for shell scripts.
package parser

import (
	"fmt"

	"github.com/afeldman/goshellcheck/internal/syntax/ast"
	"github.com/afeldman/goshellcheck/internal/syntax/lexer"
	"github.com/afeldman/goshellcheck/internal/syntax/token"
)

// Parser represents a parser for shell scripts.
type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []string
}

// New creates a new parser for the given input.
func New(input, filename string) *Parser {
	l := lexer.New(input, filename)
	p := &Parser{lexer: l}
	
	// Read two tokens to initialize currentToken and peekToken
	p.nextToken()
	p.nextToken()
	
	return p
}

// nextToken advances to the next token.
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// Parse parses the input and returns an AST.
func (p *Parser) Parse() (*ast.Program, []string) {
	program := &ast.Program{}
	
	// Parse shebang if present
	if p.currentToken.Type == token.COMMENT && len(p.currentToken.Literal) > 2 && p.currentToken.Literal[:2] == "#!" {
		program.Shebang = &ast.Shebang{
			Value:    p.currentToken.Literal[2:],
			StartPos: p.currentToken.Position,
			EndPos:   p.currentToken.Position,
		}
		p.nextToken()
	}
	
	program.StartPos = token.Position{
		Filename: p.currentToken.Position.Filename,
		Line:     1,
		Column:   1,
	}
	
	// Parse statements until EOF
	for p.currentToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}
	
	program.EndPos = p.currentToken.Position
	
	return program, p.errors
}

// parseStatement parses a statement.
func (p *Parser) parseStatement() ast.Stmt {
	switch p.currentToken.Type {
	case token.NEWLINE, token.SEMICOLON:
		p.nextToken()
		return nil
	case token.EOF:
		return nil
	default:
		return p.parseCommand()
	}
}

// parseCommand parses a command.
func (p *Parser) parseCommand() *ast.Command {
	pos := p.currentToken.Position
	
	// Try to parse a simple command
	simpleCmd := p.parseSimpleCommand()
	if simpleCmd != nil {
		return &ast.Command{
			Simple:   simpleCmd,
			StartPos: pos,
			EndPos:   simpleCmd.EndPos,
		}
	}
	
	// TODO: Parse compound commands (if, for, while, etc.)
	
	p.errors = append(p.errors, fmt.Sprintf("unexpected token: %s", p.currentToken.Literal))
	p.nextToken()
	return nil
}

// parseSimpleCommand parses a simple command.
func (p *Parser) parseSimpleCommand() *ast.SimpleCommand {
	pos := p.currentToken.Position
	
	var assignments []*ast.Assignment
	var words []*ast.Word
	
	// Parse assignments (e.g., VAR=value)
	for p.isAssignment() {
		assign := p.parseAssignment()
		if assign != nil {
			assignments = append(assignments, assign)
		}
	}
	
	// Parse command words
	for !p.isCommandTerminator() && p.currentToken.Type != token.EOF {
		word := p.parseWord()
		if word != nil {
			words = append(words, word)
		} else {
			break
		}
	}
	
	// Parse redirects
	var redirects []*ast.Redirect
	for p.isRedirect() {
		redirect := p.parseRedirect()
		if redirect != nil {
			redirects = append(redirects, redirect)
		}
	}
	
	// A simple command must have at least one word or assignment
	if len(assignments) == 0 && len(words) == 0 {
		return nil
	}
	
	endPos := pos
	if len(words) > 0 {
		endPos = words[len(words)-1].EndPos
	} else if len(assignments) > 0 {
		endPos = assignments[len(assignments)-1].EndPos
	}
	if len(redirects) > 0 {
		endPos = redirects[len(redirects)-1].EndPos
	}
	
	return &ast.SimpleCommand{
		Assignments: assignments,
		Words:       words,
		Redirects:   redirects,
		StartPos:    pos,
		EndPos:      endPos,
	}
}

// parseAssignment parses an assignment (e.g., VAR=value).
func (p *Parser) parseAssignment() *ast.Assignment {
	if !p.isAssignment() {
		return nil
	}
	
	pos := p.currentToken.Position
	
	// Parse variable name
	nameToken := p.currentToken
	if nameToken.Type != token.LITERAL {
		p.errors = append(p.errors, fmt.Sprintf("expected variable name, got %s", nameToken.Literal))
		p.nextToken()
		return nil
	}
	
	name := nameToken.Literal
	p.nextToken()
	
	// Expect '='
	if p.currentToken.Type != token.ASSIGN {
		p.errors = append(p.errors, fmt.Sprintf("expected '=', got %s", p.currentToken.Literal))
		return nil
	}
	
	p.nextToken()
	
	// Parse value
	var value *ast.Word
	if !p.isCommandTerminator() && !p.isRedirect() {
		value = p.parseWord()
	}
	
	endPos := pos.Advance(len(name))
	if value != nil {
		endPos = value.EndPos
	}
	
	return &ast.Assignment{
		Name:     name,
		Value:    value,
		StartPos: pos,
		EndPos:   endPos,
	}
}

// parseWord parses a word.
func (p *Parser) parseWord() *ast.Word {
	pos := p.currentToken.Position
	var parts []ast.WordPart
	
	for {
		var part ast.WordPart
		
		switch p.currentToken.Type {
		case token.LITERAL:
			part = &ast.Literal{
				Value:    p.currentToken.Literal,
				StartPos: p.currentToken.Position,
				EndPos:   p.currentToken.Position.Advance(len(p.currentToken.Literal)),
			}
			p.nextToken()
			
		case token.SINGLE_QUOTED:
			// Remove quotes
			value := p.currentToken.Literal
			if len(value) >= 2 && value[0] == '\'' && value[len(value)-1] == '\'' {
				value = value[1 : len(value)-1]
			}
			part = &ast.SingleQuoted{
				Value:    value,
				StartPos: p.currentToken.Position,
				EndPos:   p.currentToken.Position.Advance(len(p.currentToken.Literal)),
			}
			p.nextToken()
			
		case token.DOUBLE_QUOTED:
			// Parse double-quoted string with possible expansions
			part = p.parseDoubleQuoted()
			
		case token.DOLLAR:
			part = p.parseVariableExpansion()
			
		case token.DOLLAR_BRACE:
			part = p.parseBracedVariableExpansion()
			
		case token.DOLLAR_PAREN:
			part = p.parseCommandExpansion()
			
		default:
			// Not a word part
			if len(parts) == 0 {
				return nil
			}
			// End of word
			endPos := parts[len(parts)-1].End()
			return &ast.Word{
				Parts:    parts,
				StartPos: pos,
				EndPos:   endPos,
			}
		}
		
		if part != nil {
			parts = append(parts, part)
		}
	}
}

// parseDoubleQuoted parses a double-quoted string.
func (p *Parser) parseDoubleQuoted() *ast.DoubleQuoted {
	pos := p.currentToken.Position
	literal := p.currentToken.Literal
	
	// For now, treat the entire double-quoted string as a literal
	// TODO: Parse expansions inside double quotes
	part := &ast.Literal{
		Value:    literal,
		StartPos: pos,
		EndPos:   pos.Advance(len(literal)),
	}
	
	p.nextToken()
	
	return &ast.DoubleQuoted{
		Parts:    []ast.WordPart{part},
		StartPos: pos,
		EndPos:   pos.Advance(len(literal)),
	}
}

// parseVariableExpansion parses a simple variable expansion like $var.
func (p *Parser) parseVariableExpansion() *ast.VariableExpansion {
	pos := p.currentToken.Position
	p.nextToken() // Skip '$'
	
	if p.currentToken.Type != token.LITERAL {
		p.errors = append(p.errors, fmt.Sprintf("expected variable name after $, got %s", p.currentToken.Literal))
		return nil
	}
	
	name := p.currentToken.Literal
	endPos := p.currentToken.Position.Advance(len(name))
	p.nextToken()
	
	return &ast.VariableExpansion{
		Name:     name,
		StartPos: pos,
		EndPos:   endPos,
	}
}

// parseBracedVariableExpansion parses a braced variable expansion like ${var}.
func (p *Parser) parseBracedVariableExpansion() *ast.BracedVariableExpansion {
	pos := p.currentToken.Position
	p.nextToken() // Skip '${'
	
	// For now, parse simple ${var} without complex expansions
	if p.currentToken.Type != token.LITERAL {
		p.errors = append(p.errors, fmt.Sprintf("expected variable name after ${, got %s", p.currentToken.Literal))
		return nil
	}
	
	name := p.currentToken.Literal
	p.nextToken()
	
	// Expect '}'
	if p.currentToken.Type != token.RBRACE {
		p.errors = append(p.errors, fmt.Sprintf("expected '}', got %s", p.currentToken.Literal))
		return nil
	}
	
	endPos := p.currentToken.Position.Advance(1) // Include '}'
	p.nextToken()
	
	return &ast.BracedVariableExpansion{
		Name:     name,
		StartPos: pos,
		EndPos:   endPos,
	}
}

// parseCommandExpansion parses a command expansion like $(cmd).
func (p *Parser) parseCommandExpansion() *ast.CommandExpansion {
	pos := p.currentToken.Position
	p.nextToken() // Skip '$('
	
	// For now, parse a simple command inside $()
	// TODO: Handle nested expansions and complex commands
	var cmd *ast.Command
	if p.currentToken.Type != token.RPAREN {
		// Parse a simple command
		simpleCmd := p.parseSimpleCommand()
		if simpleCmd != nil {
			cmd = &ast.Command{
				Simple:   simpleCmd,
				StartPos: simpleCmd.StartPos,
				EndPos:   simpleCmd.EndPos,
			}
		}
	}
	
	// Expect ')'
	if p.currentToken.Type != token.RPAREN {
		p.errors = append(p.errors, fmt.Sprintf("expected ')', got %s", p.currentToken.Literal))
		return nil
	}
	
	endPos := p.currentToken.Position.Advance(1) // Include ')'
	p.nextToken()
	
	return &ast.CommandExpansion{
		Command:  cmd,
		StartPos: pos,
		EndPos:   endPos,
	}
}

// parseRedirect parses a redirect.
func (p *Parser) parseRedirect() *ast.Redirect {
	pos := p.currentToken.Position
	
	// Parse file descriptor (optional)
	fd := -1 // default
	if p.currentToken.Type == token.LITERAL {
		// Check if it's a number
		isNum := true
		for _, ch := range p.currentToken.Literal {
			if ch < '0' || ch > '9' {
				isNum = false
				break
			}
		}
		if isNum {
			// Simple conversion for now
			fd = 0 // We'll implement proper conversion later
			p.nextToken()
		}
	}
	
	// Parse redirect operator
	var operator string
	switch p.currentToken.Type {
	case token.REDIRECT_IN:
		operator = "<"
	case token.REDIRECT_OUT:
		operator = ">"
	case token.REDIRECT_APPEND:
		operator = ">>"
	case token.REDIRECT_HERE:
		operator = "<<"
	default:
		p.errors = append(p.errors, fmt.Sprintf("expected redirect operator, got %s", p.currentToken.Literal))
		return nil
	}
	
	p.nextToken()
	
	// Parse target word
	target := p.parseWord()
	if target == nil {
		p.errors = append(p.errors, "expected redirect target")
		return nil
	}
	
	endPos := target.EndPos
	return &ast.Redirect{
		FD:        fd,
		Operator:  operator,
		Target:    target,
		StartPos:  pos,
		EndPos:    endPos,
	}
}

// isAssignment checks if the current token starts an assignment.
func (p *Parser) isAssignment() bool {
	if p.currentToken.Type != token.LITERAL {
		return false
	}
	
	// Check if next token is '='
	return p.peekToken.Type == token.ASSIGN
}

// isRedirect checks if the current token is a redirect.
func (p *Parser) isRedirect() bool {
	switch p.currentToken.Type {
	case token.REDIRECT_IN, token.REDIRECT_OUT, token.REDIRECT_APPEND, token.REDIRECT_HERE:
		return true
	case token.LITERAL:
		// Could be a file descriptor number followed by redirect
		// Check if it's a number
		isNum := true
		for _, ch := range p.currentToken.Literal {
			if ch < '0' || ch > '9' {
				isNum = false
				break
			}
		}
		if isNum {
			// Check if next token is a redirect operator
			switch p.peekToken.Type {
			case token.REDIRECT_IN, token.REDIRECT_OUT, token.REDIRECT_APPEND, token.REDIRECT_HERE:
				return true
			}
		}
	}
	return false
}

// isCommandTerminator checks if the current token terminates a command.
func (p *Parser) isCommandTerminator() bool {
	switch p.currentToken.Type {
	case token.NEWLINE, token.SEMICOLON, token.PIPE, token.AND, token.OR, token.BACKGROUND, token.EOF:
		return true
	default:
		return false
	}
}
