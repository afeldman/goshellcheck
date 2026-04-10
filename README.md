# goshellcheck

A Go implementation of shell script static analysis, inspired by the original Haskell [ShellCheck](https://github.com/koalaman/shellcheck).

## Project Status

⚠️ **Early Development** - This is a work in progress. Currently only provides CLI scaffolding and basic structure.

## Installation

### From Source

```bash
git clone https://github.com/afeldman/goshellcheck
cd goshellcheck
make build
```

### Using Go Install

```bash
go install github.com/afeldman/goshellcheck/cmd/goshellcheck@latest
```

## Usage

```bash
# Check a shell script
goshellcheck script.sh

# Show version information
goshellcheck --version

# Show help
goshellcheck --help
```

## Features (Planned)

- [ ] Shell script parsing
- [ ] Common error detection
- [ ] Style suggestions
- [ ] Portability warnings
- [ ] Security checks
- [ ] Multiple output formats (JSON, CheckStyle, etc.)
- [ ] Editor integration support

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Running Tests with Coverage

```bash
make test-coverage
```

### Linting

```bash
make lint
```

### Creating a Release Snapshot

```bash
make snapshot
```

## Project Structure

```
goshellcheck/
├── cmd/goshellcheck/      # CLI entry point
├── internal/
│   ├── analyzer/         # Core analysis engine
│   ├── cli/             # Command-line interface
│   └── version/         # Version information
├── testdata/            # Test files
├── Makefile            # Build automation
└── README.md           # This file
```

## License

Apache 2.0

## Acknowledgments

- Inspired by [ShellCheck](https://github.com/koalaman/shellcheck) by Vidar Holen
- Built with Go standard library for minimal dependencies
