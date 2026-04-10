# goshellcheck Makefile
# Minimal Makefile — delegates to GoReleaser and Go directly

.PHONY: test lint clean snapshot check help

## Test — run all tests
test:
	go test ./...

## Coverage — run tests with coverage report
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

## Lint — run linters
lint:
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
		golangci-lint run ./...; \
	fi

## Check — lint + test
check: lint test

## Snapshot — build local snapshot with GoReleaser
snapshot:
	@if command -v goreleaser >/dev/null; then \
		goreleaser release --snapshot --clean; \
	else \
		echo "goreleaser not found, installing..."; \
		go install github.com/goreleaser/goreleaser/v2@latest && \
		goreleaser release --snapshot --clean; \
	fi

## Clean — remove build artifacts and temp files
clean:
	rm -rf dist/ bin/ coverage.out coverage.html

## Help — show available targets
help:
	@echo "Available targets:"
	@echo "  test       - Run all tests"
	@echo "  coverage   - Run tests with coverage HTML report"
	@echo "  lint       - Run golangci-lint"
	@echo "  check      - Run lint + test"
	@echo "  snapshot   - Build snapshot release with GoReleaser"
	@echo "  clean      - Remove build artifacts"
	@echo "  help       - Show this message"
