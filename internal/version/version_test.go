package version

import (
	"strings"
	"testing"
)

func TestInfo(t *testing.T) {
	info := Info()
	
	// Should contain version info
	if !strings.Contains(info, "goshellcheck version") {
		t.Errorf("Info() should contain 'goshellcheck version', got: %s", info)
	}
	
	// Should contain commit info
	if !strings.Contains(info, "commit:") {
		t.Errorf("Info() should contain commit info, got: %s", info)
	}
	
	// Should contain build date
	if !strings.Contains(info, "built:") {
		t.Errorf("Info() should contain build date, got: %s", info)
	}
	
	// Should contain Go version
	if !strings.Contains(info, "go:") {
		t.Errorf("Info() should contain Go version, got: %s", info)
	}
}

func TestShort(t *testing.T) {
	short := Short()
	
	// Should contain version info
	if !strings.Contains(short, "goshellcheck version") {
		t.Errorf("Short() should contain 'goshellcheck version', got: %s", short)
	}
	
	// Should be shorter than full info
	if len(short) >= len(Info()) {
		t.Errorf("Short() should be shorter than Info(), Short: %d, Info: %d", len(short), len(Info()))
	}
}
