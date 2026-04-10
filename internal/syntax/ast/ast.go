// Package ast defines the abstract syntax tree for shell scripts.
package ast

import (
	"github.com/afeldman/goshellcheck/internal/syntax/token"
)

// Node represents a node in the AST.
type Node interface {
	Pos() token.Position
	End() token.Position
}

// Expr represents an expression node.
type Expr interface {
	Node
	exprNode()
}

// Stmt represents a statement node.
type Stmt interface {
	Node
	stmtNode()
}

// Program represents a complete shell script.
type Program struct {
	Shebang   *Shebang
	Statements []Stmt
	Pos       token.Position
	EndPos    token.Position
}

func (p *Program) Pos() token.Position { return p.Pos }
func (p *Program) End() token.Position { return p.EndPos }

// Shebang represents a shebang line.
type Shebang struct {
	Value    string
	Pos      token.Position
	EndPos   token.Position
}

func (s *Shebang) Pos() token.Position { return s.Pos }
func (s *Shebang) End() token.Position { return s.EndPos }
func (s *Shebang) exprNode() {}

// Comment represents a comment.
type Comment struct {
	Text     string
	Pos      token.Position
	EndPos   token.Position
}

func (c *Comment) Pos() token.Position { return c.Pos }
func (c *Comment) End() token.Position { return c.EndPos }
func (c *Comment) exprNode() {}

// Word represents a word in a command.
type Word struct {
	Parts    []WordPart
	Pos      token.Position
	EndPos   token.Position
}

func (w *Word) Pos() token.Position { return w.Pos }
func (w *Word) End() token.Position { return w.EndPos }
func (w *Word) exprNode() {}

// WordPart represents a part of a word.
type WordPart interface {
	Node
	wordPartNode()
}

// Literal represents a literal string.
type Literal struct {
	Value    string
	Pos      token.Position
	EndPos   token.Position
}

func (l *Literal) Pos() token.Position { return l.Pos }
func (l *Literal) End() token.Position { return l.EndPos }
func (l *Literal) wordPartNode() {}
func (l *Literal) exprNode() {}

// SingleQuoted represents a single-quoted string.
type SingleQuoted struct {
	Value    string
	Pos      token.Position
	EndPos   token.Position
}

func (s *SingleQuoted) Pos() token.Position { return s.Pos }
func (s *SingleQuoted) End() token.Position { return s.EndPos }
func (s *SingleQuoted) wordPartNode() {}
func (s *SingleQuoted) exprNode() {}

// DoubleQuoted represents a double-quoted string.
type DoubleQuoted struct {
	Parts    []WordPart
	Pos      token.Position
	EndPos   token.Position
}

func (d *DoubleQuoted) Pos() token.Position { return d.Pos }
func (d *DoubleQuoted) End() token.Position { return d.EndPos }
func (d *DoubleQuoted) wordPartNode() {}
func (d *DoubleQuoted) exprNode() {}

// VariableExpansion represents a simple variable expansion like $var.
type VariableExpansion struct {
	Name     string
	Pos      token.Position
	EndPos   token.Position
}

func (v *VariableExpansion) Pos() token.Position { return v.Pos }
func (v *VariableExpansion) End() token.Position { return v.EndPos }
func (v *VariableExpansion) wordPartNode() {}
func (v *VariableExpansion) exprNode() {}

// BracedVariableExpansion represents a braced variable expansion like ${var}.
type BracedVariableExpansion struct {
	Name     string
	Content  *Word // For more complex expansions like ${var:-default}
	Pos      token.Position
	EndPos   token.Position
}

func (b *BracedVariableExpansion) Pos() token.Position { return b.Pos }
func (b *BracedVariableExpansion) End() token.Position { return b.EndPos }
func (b *BracedVariableExpansion) wordPartNode() {}
func (b *BracedVariableExpansion) exprNode() {}

// CommandExpansion represents a command expansion like $(cmd) or `cmd`.
type CommandExpansion struct {
	Command  *Command
	Pos      token.Position
	EndPos   token.Position
}

func (c *CommandExpansion) Pos() token.Position { return c.Pos }
func (c *CommandExpansion) End() token.Position { return c.EndPos }
func (c *CommandExpansion) wordPartNode() {}
func (c *CommandExpansion) exprNode() {}

// Assignment represents a variable assignment.
type Assignment struct {
	Name     string
	Value    *Word
	Pos      token.Position
	EndPos   token.Position
}

func (a *Assignment) Pos() token.Position { return a.Pos }
func (a *Assignment) End() token.Position { return a.EndPos }
func (a *Assignment) exprNode() {}

// SimpleCommand represents a simple command.
type SimpleCommand struct {
	Assignments []*Assignment
	Words       []*Word
	Redirects   []*Redirect
	Pos         token.Position
	EndPos      token.Position
}

func (s *SimpleCommand) Pos() token.Position { return s.Pos }
func (s *SimpleCommand) End() token.Position { return s.EndPos }
func (s *SimpleCommand) stmtNode() {}

// Pipeline represents a pipeline of commands.
type Pipeline struct {
	Commands []*Command
	Negated  bool // true if preceded by !
	Pos      token.Position
	EndPos   token.Position
}

func (p *Pipeline) Pos() token.Position { return p.Pos }
func (p *Pipeline) End() token.Position { return p.EndPos }
func (p *Pipeline) stmtNode() {}

// Command represents a command (simple command or compound command).
type Command struct {
	Simple *SimpleCommand
	// Future: compound commands like if, for, while, etc.
	Pos    token.Position
	EndPos token.Position
}

func (c *Command) Pos() token.Position { return c.Pos }
func (c *Command) End() token.Position { return c.EndPos }
func (c *Command) stmtNode() {}

// Redirect represents an I/O redirection.
type Redirect struct {
	FD       int    // file descriptor, -1 for default
	Operator string // <, >, >>, <<, etc.
	Target   *Word
	Pos      token.Position
	EndPos   token.Position
}

func (r *Redirect) Pos() token.Position { return r.Pos }
func (r *Redirect) End() token.Position { return r.EndPos }
func (r *Redirect) exprNode() {}

// CommandList represents a list of commands separated by ; or newline.
type CommandList struct {
	Commands []Stmt
	Pos      token.Position
	EndPos   token.Position
}

func (c *CommandList) Pos() token.Position { return c.Pos }
func (c *CommandList) End() token.Position { return c.EndPos }
func (c *CommandList) stmtNode() {}

// AndOr represents an AND-OR list (commands separated by && or ||).
type AndOr struct {
	Left     Stmt
	Operator string // && or ||
	Right    Stmt
	Pos      token.Position
	EndPos   token.Position
}

func (a *AndOr) Pos() token.Position { return a.Pos }
func (a *AndOr) End() token.Position { return a.EndPos }
func (a *AndOr) stmtNode() {}
