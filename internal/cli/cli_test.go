package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantHelp    bool
		wantVersion bool
		wantFiles   []string
		wantErr     bool
	}{
		{
			name:      "no arguments",
			args:      []string{},
			wantHelp:  false,
			wantFiles: []string{},
		},
		{
			name:      "help flag",
			args:      []string{"--help"},
			wantHelp:  true,
			wantFiles: []string{},
		},
		{
			name:        "version flag",
			args:        []string{"--version"},
			wantVersion: true,
			wantFiles:   []string{},
		},
		{
			name:      "short help flag",
			args:      []string{"-h"},
			wantHelp:  true,
			wantFiles: []string{},
		},
		{
			name:        "short version flag",
			args:        []string{"-v"},
			wantVersion: true,
			wantFiles:   []string{},
		},
		{
			name:      "single file",
			args:      []string{"script.sh"},
			wantFiles: []string{"script.sh"},
		},
		{
			name:      "multiple files",
			args:      []string{"script1.sh", "script2.sh"},
			wantFiles: []string{"script1.sh", "script2.sh"},
		},
		{
			name:      "flags and files",
			args:      []string{"--help", "script.sh"},
			wantHelp:  true,
			wantFiles: []string{"script.sh"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := Parse(tt.args)
			
			if tt.wantErr && err == nil {
				t.Errorf("Parse() expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
			}
			if err != nil {
				return
			}
			
			if cfg.Help != tt.wantHelp {
				t.Errorf("Parse() Help = %v, want %v", cfg.Help, tt.wantHelp)
			}
			if cfg.Version != tt.wantVersion {
				t.Errorf("Parse() Version = %v, want %v", cfg.Version, tt.wantVersion)
			}
			if len(cfg.Files) != len(tt.wantFiles) {
				t.Errorf("Parse() Files length = %d, want %d", len(cfg.Files), len(tt.wantFiles))
			}
			for i, file := range cfg.Files {
				if file != tt.wantFiles[i] {
					t.Errorf("Parse() Files[%d] = %s, want %s", i, file, tt.wantFiles[i])
				}
			}
		})
	}
}

func TestRunHelp(t *testing.T) {
	// Capture stderr (flag.Usage() writes to stderr)
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	
	cfg := &Config{Help: true}
	exitCode := Run(cfg)
	
	w.Close()
	os.Stderr = oldStderr
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	
	if exitCode != 0 {
		t.Errorf("Run() with Help=true should return 0, got %d", exitCode)
	}
	
	if !strings.Contains(output, "Usage:") {
		t.Errorf("Run() with Help=true should output usage, got: %s", output)
	}
}

func TestRunVersion(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	cfg := &Config{Version: true}
	exitCode := Run(cfg)
	
	w.Close()
	os.Stdout = oldStdout
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	
	if exitCode != 0 {
		t.Errorf("Run() with Version=true should return 0, got %d", exitCode)
	}
	
	if !strings.Contains(output, "goshellcheck version") {
		t.Errorf("Run() with Version=true should output version info, got: %s", output)
	}
}

func TestRunNoFiles(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	
	cfg := &Config{Files: []string{}}
	exitCode := Run(cfg)
	
	w.Close()
	os.Stderr = oldStderr
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	
	if exitCode != 2 {
		t.Errorf("Run() with no files should return 2, got %d", exitCode)
	}
	
	if !strings.Contains(output, "Error: no files specified") {
		t.Errorf("Run() with no files should output error, got: %s", output)
	}
}

func TestRunFileDoesNotExist(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()
	nonExistentFile := filepath.Join(tmpDir, "nonexistent.sh")
	
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	
	cfg := &Config{Files: []string{nonExistentFile}}
	exitCode := Run(cfg)
	
	w.Close()
	os.Stderr = oldStderr
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	
	if exitCode != 1 {
		t.Errorf("Run() with non-existent file should return 1, got %d", exitCode)
	}
	
	expectedError := fmt.Sprintf("Error: file does not exist: %s", nonExistentFile)
	if !strings.Contains(output, expectedError) {
		t.Errorf("Run() should output file not found error, got: %s", output)
	}
}

func TestRunValidFiles(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir := t.TempDir()
	
	file1 := filepath.Join(tmpDir, "test1.sh")
	file2 := filepath.Join(tmpDir, "test2.sh")
	
	// Create test files
	os.WriteFile(file1, []byte("#!/bin/bash\necho 'test1'"), 0644)
	os.WriteFile(file2, []byte("#!/bin/sh\necho 'test2'"), 0644)
	
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	cfg := &Config{Files: []string{file1, file2}}
	exitCode := Run(cfg)
	
	w.Close()
	os.Stdout = oldStdout
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	
	if exitCode != 0 {
		t.Errorf("Run() with valid files should return 0, got %d", exitCode)
	}
	
	if !strings.Contains(output, "All files passed analysis") {
		t.Errorf("Run() should indicate all files passed, got: %s", output)
	}
}
