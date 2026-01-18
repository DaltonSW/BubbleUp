.PHONY: help test test-all lint build clean fmt vet tidy ensure-valid ci

LINTER = "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2"

# Default target
help:
	@echo "Available targets:"
	@echo "  make test         - Run unit tests"
	@echo "  make test-all     - Run all tests"
	@echo "  make lint         - Run golangci-lint"
	@echo "  make fmt          - Format code with gofmt"
	@echo "  make vet          - Run go vet"
	@echo "  make tidy         - Run go mod tidy"
	@echo "  make build        - Build the package"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make ensure-valid - Run tidy, test, lint, and vet"
	@echo "  make ci           - Run all CI checks (fmt, vet, lint, test-all)"

ensure-valid: tidy test lint vet

# Run unit tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=./coverage.txt -covermode=atomic ./...

# Run all tests
test-all: test

# Run linter
lint:
	@echo "Running linter..."
	@go run $(LINTER) run ./... --timeout=5m

# Format code
fmt:
	@echo "Formatting code..."
	@gofmt -s -w .

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run go mod tidy
tidy:
	@echo "Running go mod tidy..."
	@go mod tidy

# Build the package
build:
	@echo "Building package..."
	@go build ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@go clean
	@rm -f coverage.txt

# Run all CI checks locally
ci: fmt vet lint test-all
	@echo "All CI checks passed!"
