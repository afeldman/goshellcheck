# goshellcheck Makefile

.PHONY: all build test test-coverage clean lint snapshot help

# Build variables
BINARY_NAME := goshellcheck
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -X 'github.com/afeldman/goshellcheck/internal/version.Version=$(VERSION)' \
           -X 'github.com/afeldman/goshellcheck/internal/version.Commit=$(COMMIT)' \
           -X 'github.com/afeldman/goshellcheck/internal/version.Date=$(DATE)'

all: build

## Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/goshellcheck

## Install the binary to $GOPATH/bin
install:
	go install -ldflags "$(LDFLAGS)" ./cmd/goshellcheck

## Run tests
test:
	go test ./...

## Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html dist/

## Run linters (requires golangci-lint)
lint:
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run ./...; \
	fi

## Create a local snapshot release with GoReleaser
snapshot:
	@if command -v goreleaser >/dev/null; then \
		goreleaser release --snapshot --clean; \
	else \
		echo "goreleaser not installed. Installing..."; \
		go install github.com/goreleaser/goreleaser/v2@latest; \
		goreleaser release --snapshot --clean; \
	fi

## Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  install        - Install the binary to \$GOPATH/bin"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Clean build artifacts"
	@echo "  lint           - Run linters (requires golangci-lint)"
	@echo "  snapshot       - Create a local snapshot release"
	@echo "  help           - Show this help message"
