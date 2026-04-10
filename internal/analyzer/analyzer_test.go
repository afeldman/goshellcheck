package analyzer

import (
	"os"
	"strings"
	"testing"
)

func TestNewAnalyzer(t *testing.T) {
	analyzer := New()
	if analyzer == nil {
		t.Error("New() should return a non-nil Analyzer")
	}
}

func TestAnalyzeReader(t *testing.T) {
	analyzer := New()
	
	// Test with empty content
	reader := strings.NewReader("")
	result, err := analyzer.AnalyzeReader("test.sh", reader)
	
	if err != nil {
		t.Errorf("AnalyzeReader() unexpected error: %v", err)
	}
	
	if result.FilePath != "test.sh" {
		t.Errorf("AnalyzeReader() FilePath = %s, want test.sh", result.FilePath)
	}
	
	if len(result.Errors) != 0 {
		t.Errorf("AnalyzeReader() should return no errors for empty content, got %d", len(result.Errors))
	}
	
	if len(result.Warnings) != 0 {
		t.Errorf("AnalyzeReader() should return no warnings for empty content, got %d", len(result.Warnings))
	}
}

func TestAnalyzeFile(t *testing.T) {
	analyzer := New()
	
	// Create a temporary file
	tmpFile := t.TempDir() + "/test.sh"
	content := "#!/bin/bash\necho 'Hello, World!'"
	
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	result, err := analyzer.AnalyzeFile(tmpFile)
	
	if err != nil {
		t.Errorf("AnalyzeFile() unexpected error: %v", err)
	}
	
	if result.FilePath != tmpFile {
		t.Errorf("AnalyzeFile() FilePath = %s, want %s", result.FilePath, tmpFile)
	}
	
	// Should return empty results for now (placeholder implementation)
	if len(result.Errors) != 0 {
		t.Errorf("AnalyzeFile() should return no errors for placeholder, got %d", len(result.Errors))
	}
	
	if len(result.Warnings) != 0 {
		t.Errorf("AnalyzeFile() should return no warnings for placeholder, got %d", len(result.Warnings))
	}
}

func TestAnalyzeFileNotFound(t *testing.T) {
	analyzer := New()
	
	// Try to analyze non-existent file
	nonExistentFile := t.TempDir() + "/nonexistent.sh"
	_, err := analyzer.AnalyzeFile(nonExistentFile)
	
	if err == nil {
		t.Error("AnalyzeFile() should return error for non-existent file")
	}
	
	if !strings.Contains(err.Error(), "failed to open file") {
		t.Errorf("AnalyzeFile() error should mention 'failed to open file', got: %v", err)
	}
}
