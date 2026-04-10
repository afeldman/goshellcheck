// Package analyzer provides the core analysis pipeline for goshellcheck.
package analyzer

import (
	"fmt"
	"io"
	"os"

	"github.com/afeldman/goshellcheck/internal/diag"
	"github.com/afeldman/goshellcheck/internal/lint"
	"github.com/afeldman/goshellcheck/internal/lint/rules"
	"github.com/afeldman/goshellcheck/internal/syntax/parser"
)

// Analyzer is the main analysis engine.
type Analyzer struct {
	lintEngine *lint.Engine
}

// New creates a new Analyzer with default rules.
func New() *Analyzer {
	engine := lint.NewEngine()
	for _, rule := range rules.All() {
		engine.AddRule(rule)
	}
	return &Analyzer{lintEngine: engine}
}

// AnalyzeFile analyzes a single file.
func (a *Analyzer) AnalyzeFile(filePath string) ([]diag.Diagnostic, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	return a.AnalyzeReader(filePath, file)
}

// AnalyzeReader analyzes content from a reader.
func (a *Analyzer) AnalyzeReader(filePath string, r io.Reader) ([]diag.Diagnostic, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read content: %w", err)
	}

	// Parse the shell script
	p := parser.New(string(content), filePath)
	prog, parseErrors := p.Parse()
	
	var diagnostics []diag.Diagnostic
	
	// Add parse errors as diagnostics
	for _, errMsg := range parseErrors {
		diagnostics = append(diagnostics, diag.New(
			"SC1000",
			diag.Error,
			errMsg,
			filePath,
			1, // Default line
			1, // Default column
		))
	}
	
	// Run lint rules if parsing was successful
	if prog != nil {
		ruleDiags := a.lintEngine.Analyze(prog, filePath)
		diagnostics = append(diagnostics, ruleDiags...)
	}

	return diagnostics, nil
}
