// Package cli provides command-line interface functionality for goshellcheck.
package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/afeldman/goshellcheck/internal/analyzer"
	"github.com/afeldman/goshellcheck/internal/version"
)

// Config holds the parsed command-line configuration.
type Config struct {
	Help    bool
	Version bool
	Files   []string
}

// Parse parses command-line arguments and returns a Config.
func Parse(args []string) (*Config, error) {
	cfg := &Config{}

	fs := flag.NewFlagSet("goshellcheck", flag.ContinueOnError)
	fs.BoolVar(&cfg.Help, "help", false, "Show help")
	fs.BoolVar(&cfg.Help, "h", false, "Show help (shorthand)")
	fs.BoolVar(&cfg.Version, "version", false, "Show version information")
	fs.BoolVar(&cfg.Version, "v", false, "Show version information (shorthand)")

	// Custom usage function
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] FILE...\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s script.sh\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "  %s --version\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "  %s --help\n", filepath.Base(os.Args[0]))
	}

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg.Files = fs.Args()

	return cfg, nil
}

// Run executes the CLI based on the provided configuration.
func Run(cfg *Config) int {
	if cfg.Help {
		printUsage()
		return 0
	}

	if cfg.Version {
		fmt.Println(version.Info())
		return 0
	}

	if len(cfg.Files) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no files specified\n")
		printUsage()
		return 2
	}

	// Check if files exist
	hasErrors := false
	for _, file := range cfg.Files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: file does not exist: %s\n", file)
			hasErrors = true
		}
	}
	if hasErrors {
		return 1
	}

	// Analyze files
	analyzer := analyzer.New()
	allPassed := true

	for _, file := range cfg.Files {
		result, err := analyzer.AnalyzeFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error analyzing %s: %v\n", file, err)
			allPassed = false
			continue
		}

		// Print results (placeholder - will be formatted better later)
		if len(result.Errors) > 0 || len(result.Warnings) > 0 {
			allPassed = false
			fmt.Printf("\n%s:\n", file)
			for _, err := range result.Errors {
				fmt.Printf("  Error: %s (line %d, col %d)\n", err.Message, err.Line, err.Column)
			}
			for _, warn := range result.Warnings {
				fmt.Printf("  Warning: %s (line %d, col %d)\n", warn.Message, warn.Line, warn.Column)
			}
		}
	}

	if allPassed {
		fmt.Println("All files passed analysis.")
		return 0
	}
	return 1
}

// printUsage prints the usage information to stderr.
func printUsage() {
	fs := flag.NewFlagSet("goshellcheck", flag.ContinueOnError)
	fs.Bool("help", false, "Show help")
	fs.Bool("h", false, "Show help (shorthand)")
	fs.Bool("version", false, "Show version information")
	fs.Bool("v", false, "Show version information (shorthand)")
	
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] FILE...\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	fs.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  %s script.sh\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "  %s --version\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "  %s --help\n", filepath.Base(os.Args[0]))
}
