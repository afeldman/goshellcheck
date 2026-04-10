// Package rules contains linting rules for shell scripts.
package rules

import (
	"github.com/afeldman/goshellcheck/internal/diag"
	"github.com/afeldman/goshellcheck/internal/lint"
	"github.com/afeldman/goshellcheck/internal/syntax/ast"
	"strings"
)

// All returns all available rules.
func All() []lint.Rule {
	return []lint.Rule{
		NewUnquotedVariableExpansionRule(),
		NewBackticksRule(),
		NewUselessCatRule(),
		NewUnquotedEchoVariableRule(),
		NewLsGrepAntiPatternRule(),
	}
}

// UnquotedVariableExpansionRule checks for unquoted variable expansions in command arguments.
type UnquotedVariableExpansionRule struct{}

func NewUnquotedVariableExpansionRule() lint.Rule {
	return lint.NewRule(
		"SC2086",
		"Double quote to prevent globbing and word splitting",
		func(prog *ast.Program, file string) []diag.Diagnostic {
			var diags []diag.Diagnostic
			
			visitor := &lint.Visitor{
				VisitNode: func(node ast.Node) bool {
					if word, ok := node.(*ast.Word); ok {
						// Check if this word is in a command argument (not an assignment)
						// For now, we'll check all words
						if hasUnquotedVariableExpansion(word) {
							pos := word.Pos()
							diags = append(diags, diag.New(
								"SC2086",
								diag.Info,
								"Double quote to prevent globbing and word splitting",
								file,
								pos.Line,
								pos.Column,
							).WithSuggestion(`Use "$var" instead of $var`))
						}
					}
					return true
				},
			}
			
			visitor.Walk(prog)
			return diags
		},
	)
}

func hasUnquotedVariableExpansion(word *ast.Word) bool {
	for _, part := range word.Parts {
		switch part.(type) {
		case *ast.VariableExpansion, *ast.BracedVariableExpansion:
			// Check if the word is quoted
			if len(word.Parts) == 1 {
				// Single variable expansion - should be quoted
				return true
			}
			// Multiple parts - check if any are not quoted
			return true
		}
	}
	return false
}

// BackticksRule checks for backticks and suggests $(...) instead.
type BackticksRule struct{}

func NewBackticksRule() lint.Rule {
	return lint.NewRule(
		"SC2006",
		"Use $(...) notation instead of legacy backticks",
		func(prog *ast.Program, file string) []diag.Diagnostic {
			var diags []diag.Diagnostic
			
			visitor := &lint.Visitor{
				VisitNode: func(node ast.Node) bool {
					// Note: The current AST doesn't distinguish between $(...) and `...`
					// Both are represented as CommandExpansion.
					// We would need to extend the parser to track which syntax was used.
					// For now, this is a placeholder.
					return true
				},
			}
			
			visitor.Walk(prog)
			return diags
		},
	)
}

// UselessCatRule checks for useless use of cat in pipelines.
type UselessCatRule struct{}

func NewUselessCatRule() lint.Rule {
	return lint.NewRule(
		"SC2002",
		"Useless cat. Consider using redirection or command arguments",
		func(prog *ast.Program, file string) []diag.Diagnostic {
			var diags []diag.Diagnostic
			
			visitor := &lint.Visitor{
				VisitNode: func(node ast.Node) bool {
					if pipeline, ok := node.(*ast.Pipeline); ok && len(pipeline.Commands) >= 2 {
						// Check if first command is "cat" with a single file argument
						if firstCmd := pipeline.Commands[0]; firstCmd.Simple != nil {
							if isCatCommand(firstCmd.Simple) {
								pos := pipeline.Pos()
								diags = append(diags, diag.New(
									"SC2002",
									diag.Style,
									"Useless cat. Consider using redirection or command arguments",
									file,
									pos.Line,
									pos.Column,
								).WithSuggestion(`Use "< file cmd" instead of "cat file | cmd"`))
							}
						}
					}
					return true
				},
			}
			
			visitor.Walk(prog)
			return diags
		},
	)
}

func isCatCommand(cmd *ast.SimpleCommand) bool {
	if len(cmd.Words) == 0 {
		return false
	}
	
	// Get the command name
	firstWord := cmd.Words[0]
	for _, part := range firstWord.Parts {
		if lit, ok := part.(*ast.Literal); ok && lit.Value == "cat" {
			// cat with a single file argument (not options)
			if len(cmd.Words) == 2 {
				// Could check if second word looks like a filename
				return true
			}
		}
	}
	return false
}

// UnquotedEchoVariableRule checks for echo $var without quotes.
type UnquotedEchoVariableRule struct{}

func NewUnquotedEchoVariableRule() lint.Rule {
	return lint.NewRule(
		"SC2028",
		"echo may not expand escape sequences. Use printf",
		func(prog *ast.Program, file string) []diag.Diagnostic {
			var diags []diag.Diagnostic
			
			visitor := &lint.Visitor{
				VisitNode: func(node ast.Node) bool {
					if cmd, ok := node.(*ast.SimpleCommand); ok && len(cmd.Words) > 0 {
						// Check if command is "echo"
						if isEchoCommand(cmd) {
							// Check arguments for unquoted variables
							for _, word := range cmd.Words[1:] {
								if hasUnquotedVariableExpansion(word) {
									pos := word.Pos()
									diags = append(diags, diag.New(
										"SC2028",
										diag.Info,
										"echo may not expand escape sequences. Use printf for complex output",
										file,
										pos.Line,
										pos.Column,
									).WithSuggestion(`Use printf or quote variables: echo "$var"`))
									break // Only report once per echo command
								}
							}
						}
					}
					return true
				},
			}
			
			visitor.Walk(prog)
			return diags
		},
	)
}

func isEchoCommand(cmd *ast.SimpleCommand) bool {
	if len(cmd.Words) == 0 {
		return false
	}
	
	firstWord := cmd.Words[0]
	for _, part := range firstWord.Parts {
		if lit, ok := part.(*ast.Literal); ok && lit.Value == "echo" {
			return true
		}
	}
	return false
}

// LsGrepAntiPatternRule checks for ls | grep patterns when globbing could be used.
type LsGrepAntiPatternRule struct{}

func NewLsGrepAntiPatternRule() lint.Rule {
	return lint.NewRule(
		"SC2010",
		"Don't use ls | grep. Use a glob or a loop",
		func(prog *ast.Program, file string) []diag.Diagnostic {
			var diags []diag.Diagnostic
			
			visitor := &lint.Visitor{
				VisitNode: func(node ast.Node) bool {
					if pipeline, ok := node.(*ast.Pipeline); ok && len(pipeline.Commands) == 2 {
						// Check if first command is "ls" and second is "grep"
						if firstCmd := pipeline.Commands[0]; firstCmd.Simple != nil {
							if secondCmd := pipeline.Commands[1]; secondCmd.Simple != nil {
								if isLsCommand(firstCmd.Simple) && isGrepCommand(secondCmd.Simple) {
									pos := pipeline.Pos()
									diags = append(diags, diag.New(
										"SC2010",
										diag.Style,
										"Don't use ls | grep. Use a glob or a loop",
										file,
										pos.Line,
										pos.Column,
									).WithSuggestion(`Use "*.ext" or find instead of "ls | grep .ext"`))
								}
							}
						}
					}
					return true
				},
			}
			
			visitor.Walk(prog)
			return diags
		},
	)
}

func isLsCommand(cmd *ast.SimpleCommand) bool {
	if len(cmd.Words) == 0 {
		return false
	}
	
	firstWord := cmd.Words[0]
	for _, part := range firstWord.Parts {
		if lit, ok := part.(*ast.Literal); ok && lit.Value == "ls" {
			return true
		}
	}
	return false
}

func isGrepCommand(cmd *ast.SimpleCommand) bool {
	if len(cmd.Words) == 0 {
		return false
	}
	
	firstWord := cmd.Words[0]
	for _, part := range firstWord.Parts {
		if lit, ok := part.(*ast.Literal); ok && strings.HasPrefix(lit.Value, "grep") {
			return true
		}
	}
	return false
}
