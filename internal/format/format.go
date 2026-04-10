// Package format provides output formatters for diagnostics.
package format

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/afeldman/goshellcheck/internal/diag"
)

// Formatter is an interface for formatting diagnostics.
type Formatter interface {
	// Format writes the diagnostics to the writer.
	Format(w io.Writer, diagnostics []diag.Diagnostic) error
}

// TTYFormatter formats diagnostics for terminal output with colors.
type TTYFormatter struct {
	// Colorize enables ANSI color codes (default: true if output is a terminal).
	Colorize bool
}

// Format implements the Formatter interface.
func (f *TTYFormatter) Format(w io.Writer, diagnostics []diag.Diagnostic) error {
	for _, d := range diagnostics {
		// Format: filename:line:col: severity: message [code]
		line := fmt.Sprintf("%s:%d:%d: %s: %s [%s]",
			d.File, d.Line, d.Column,
			d.Severity.String(), d.Message, d.Code)
		
		if d.Suggestion != "" {
			line += fmt.Sprintf(" (%s)", d.Suggestion)
		}
		
		fmt.Fprintln(w, line)
	}
	return nil
}

// GCCFormatter formats diagnostics in GCC-style output.
type GCCFormatter struct{}

// Format implements the Formatter interface.
func (f *GCCFormatter) Format(w io.Writer, diagnostics []diag.Diagnostic) error {
	for _, d := range diagnostics {
		// GCC format: filename:line:col: severity: message
		line := fmt.Sprintf("%s:%d:%d: %s: %s [%s]",
			d.File, d.Line, d.Column,
			d.Severity.String(), d.Message, d.Code)
		
		fmt.Fprintln(w, line)
	}
	return nil
}

// JSONDiagnostic is the JSON representation of a diagnostic.
type JSONDiagnostic struct {
	Code      string `json:"code"`
	Severity  string `json:"level"`
	Message   string `json:"message"`
	File      string `json:"file"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	EndLine   int    `json:"endLine,omitempty"`
	EndColumn int    `json:"endColumn,omitempty"`
	Suggestion string `json:"fix,omitempty"`
}

// JSONFormatter formats diagnostics as JSON.
type JSONFormatter struct {
	// Pretty enables pretty-printed JSON (default: false).
	Pretty bool
}

// Format implements the Formatter interface.
func (f *JSONFormatter) Format(w io.Writer, diagnostics []diag.Diagnostic) error {
	jsonDiags := make([]JSONDiagnostic, len(diagnostics))
	for i, d := range diagnostics {
		jsonDiags[i] = JSONDiagnostic{
			Code:      d.Code,
			Severity:  d.Severity.String(),
			Message:   d.Message,
			File:      d.File,
			Line:      d.Line,
			Column:    d.Column,
			EndLine:   d.EndLine,
			EndColumn: d.EndColumn,
			Suggestion: d.Suggestion,
		}
	}
	
	var encoder *json.Encoder
	if f.Pretty {
		encoder = json.NewEncoder(w)
		encoder.SetIndent("", "  ")
	} else {
		encoder = json.NewEncoder(w)
	}
	
	return encoder.Encode(jsonDiags)
}

// NewFormatter creates a formatter based on the format name.
func NewFormatter(format string) (Formatter, error) {
	switch strings.ToLower(format) {
	case "tty":
		return &TTYFormatter{}, nil
	case "gcc":
		return &GCCFormatter{}, nil
	case "json":
		return &JSONFormatter{}, nil
	default:
		return nil, fmt.Errorf("unknown format: %s", format)
	}
}
