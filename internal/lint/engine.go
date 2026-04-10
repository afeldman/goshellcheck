// Package lint provides the rule engine for shell script analysis.
package lint

import (
	"github.com/afeldman/goshellcheck/internal/diag"
	"github.com/afeldman/goshellcheck/internal/syntax/ast"
)

// Rule represents a linting rule that can analyze a shell script AST.
type Rule interface {
	// Name returns the unique name of the rule.
	Name() string
	// Description returns a brief description of what the rule checks.
	Description() string
	// Check analyzes the AST and returns a list of diagnostics.
	Check(prog *ast.Program, file string) []diag.Diagnostic
}

// Engine is the main linting engine that runs multiple rules.
type Engine struct {
	rules []Rule
}

// NewEngine creates a new linting engine with the given rules.
func NewEngine(rules ...Rule) *Engine {
	return &Engine{rules: rules}
}

// AddRule adds a rule to the engine.
func (e *Engine) AddRule(rule Rule) {
	e.rules = append(e.rules, rule)
}

// Analyze runs all rules on the given AST and returns all diagnostics.
func (e *Engine) Analyze(prog *ast.Program, file string) []diag.Diagnostic {
	var diagnostics []diag.Diagnostic
	for _, rule := range e.rules {
		diagnostics = append(diagnostics, rule.Check(prog, file)...)
	}
	return diagnostics
}

// RuleFunc is a function type that implements the Rule interface.
type RuleFunc struct {
	name        string
	description string
	check       func(prog *ast.Program, file string) []diag.Diagnostic
}

// NewRule creates a new rule from a function.
func NewRule(name, description string, check func(prog *ast.Program, file string) []diag.Diagnostic) Rule {
	return &RuleFunc{
		name:        name,
		description: description,
		check:       check,
	}
}

// Name implements the Rule interface.
func (r *RuleFunc) Name() string {
	return r.name
}

// Description implements the Rule interface.
func (r *RuleFunc) Description() string {
	return r.description
}

// Check implements the Rule interface.
func (r *RuleFunc) Check(prog *ast.Program, file string) []diag.Diagnostic {
	return r.check(prog, file)
}

// Visitor is a helper for traversing the AST.
type Visitor struct {
	// VisitNode is called for each node in the AST.
	// If it returns true, the visitor will continue to children.
	VisitNode func(node ast.Node) bool
}

// Walk traverses the AST depth-first, calling VisitNode for each node.
func (v *Visitor) Walk(node ast.Node) {
	if node == nil {
		return
	}
	
	if v.VisitNode != nil && !v.VisitNode(node) {
		return
	}
	
	switch n := node.(type) {
	case *ast.Program:
		for _, stmt := range n.Statements {
			v.Walk(stmt)
		}
	
	case *ast.SimpleCommand:
		for _, word := range n.Words {
			v.Walk(word)
		}
		for _, assign := range n.Assignments {
			v.Walk(assign)
		}
		for _, redirect := range n.Redirects {
			v.Walk(redirect)
		}
	
	case *ast.Pipeline:
		for _, cmd := range n.Commands {
			v.Walk(cmd)
		}
	
	case *ast.Command:
		if n.Simple != nil {
			v.Walk(n.Simple)
		}
	
	case *ast.Word:
		for _, part := range n.Parts {
			v.Walk(part)
		}
	
	case *ast.Assignment:
		if n.Value != nil {
			v.Walk(n.Value)
		}
	
	case *ast.Redirect:
		if n.Target != nil {
			v.Walk(n.Target)
		}
	
	case *ast.DoubleQuoted:
		for _, part := range n.Parts {
			v.Walk(part)
		}
	
	case *ast.CommandExpansion:
		if n.Command != nil {
			v.Walk(n.Command)
		}
	
	case *ast.CommandList:
		for _, cmd := range n.Commands {
			v.Walk(cmd)
		}
	
	case *ast.AndOr:
		v.Walk(n.Left)
		v.Walk(n.Right)
	
	// Leaf nodes that don't have children
	case *ast.Literal, *ast.SingleQuoted, *ast.VariableExpansion,
		*ast.BracedVariableExpansion, *ast.Comment, *ast.Shebang:
		// No children to visit
	}
}
