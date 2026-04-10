// Package diag defines diagnostic types for shell script analysis.
package diag

import (
	"fmt"
	"strings"
)

// Severity represents the severity level of a diagnostic.
type Severity int

const (
	// Error indicates a serious issue that likely causes incorrect behavior.
	Error Severity = iota
	// Warning indicates a potential issue or bad practice.
	Warning
	// Info indicates a suggestion or style issue.
	Info
	// Style indicates a purely stylistic issue.
	Style
)

// String returns the string representation of the severity.
func (s Severity) String() string {
	switch s {
	case Error:
		return "error"
	case Warning:
		return "warning"
	case Info:
		return "info"
	case Style:
		return "style"
	default:
		return "unknown"
	}
}

// Diagnostic represents a single diagnostic message.
type Diagnostic struct {
	// Code is a unique identifier for this type of diagnostic (e.g., "SC1000").
	Code string
	// Severity is the severity level of the diagnostic.
	Severity Severity
	// Message is the human-readable description of the issue.
	Message string
	// File is the path to the file where the issue was found.
	File string
	// Line is the 1-based line number where the issue was found.
	Line int
	// Column is the 1-based column number where the issue was found.
	Column int
	// EndLine is the 1-based line number where the issue ends (optional).
	EndLine int
	// EndColumn is the 1-based column number where the issue ends (optional).
	EndColumn int
	// Suggestion is an optional suggested fix.
	Suggestion string
}

// New creates a new diagnostic.
func New(code string, severity Severity, message, file string, line, column int) Diagnostic {
	return Diagnostic{
		Code:     code,
		Severity: severity,
		Message:  message,
		File:     file,
		Line:     line,
		Column:   column,
	}
}

// WithEnd adds end position to the diagnostic.
func (d Diagnostic) WithEnd(endLine, endColumn int) Diagnostic {
	d.EndLine = endLine
	d.EndColumn = endColumn
	return d
}

// WithSuggestion adds a suggestion to the diagnostic.
func (d Diagnostic) WithSuggestion(suggestion string) Diagnostic {
	d.Suggestion = suggestion
	return d
}

// HasPosition returns true if the diagnostic has position information.
func (d Diagnostic) HasPosition() bool {
	return d.Line > 0 && d.Column > 0
}

// HasEndPosition returns true if the diagnostic has end position information.
func (d Diagnostic) HasEndPosition() bool {
	return d.EndLine > 0 && d.EndColumn > 0
}

// String returns a simple string representation of the diagnostic.
func (d Diagnostic) String() string {
	var parts []string
	if d.File != "" {
		parts = append(parts, d.File)
	}
	if d.HasPosition() {
		parts = append(parts, fmt.Sprintf("%d:%d", d.Line, d.Column))
	}
	if d.Code != "" {
		parts = append(parts, fmt.Sprintf("[%s]", d.Code))
	}
	parts = append(parts, d.Severity.String()+":", d.Message)
	if d.Suggestion != "" {
		parts = append(parts, fmt.Sprintf("(%s)", d.Suggestion))
	}
	return strings.Join(parts, " ")
}

// DiagnosticList is a list of diagnostics.
type DiagnosticList []Diagnostic

// FilterBySeverity returns a new list containing only diagnostics with the given severity or higher.
func (dl DiagnosticList) FilterBySeverity(minSeverity Severity) DiagnosticList {
	var result DiagnosticList
	for _, d := range dl {
		if d.Severity <= minSeverity {
			result = append(result, d)
		}
	}
	return result
}

// HasErrors returns true if the list contains any errors.
func (dl DiagnosticList) HasErrors() bool {
	for _, d := range dl {
		if d.Severity == Error {
			return true
		}
	}
	return false
}

// HasWarningsOrErrors returns true if the list contains any warnings or errors.
func (dl DiagnosticList) HasWarningsOrErrors() bool {
	for _, d := range dl {
		if d.Severity <= Warning {
			return true
		}
	}
	return false
}
