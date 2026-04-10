// Package version provides version information for goshellcheck.
package version

import (
	"fmt"
	"runtime"
)

// These variables are set via ldflags during build.
var (
	Version   = "dev"
	Commit    = "none"
	Date      = "unknown"
	GoVersion = runtime.Version()
)

// Info returns version information as a string.
func Info() string {
	return fmt.Sprintf("goshellcheck version %s (commit: %s, built: %s, go: %s)",
		Version, Commit, Date, GoVersion)
}

// Short returns a short version string.
func Short() string {
	return fmt.Sprintf("goshellcheck version %s", Version)
}
