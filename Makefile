.PHONY: build test clean install lint fmt vet

# Build the CLI tool
build:
	go build -o bin/toonify ./cmd/toonify

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install the CLI tool
install:
	go install ./cmd/toonify

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run all checks
check: fmt vet test

# Build example
example:
	go run examples/basic/main.go

# Test CLI with sample data
test-cli: build
	./bin/toonify -input testdata/sample.json
	echo '{"name":"Alice","age":30}' | ./bin/toonify

# Initialize go modules
mod-init:
	go mod init github.com/Palaciodiego008/toonify

# Tidy go modules
mod-tidy:
	go mod tidy

# Download dependencies
mod-download:
	go mod download

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the CLI tool"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  install       - Install the CLI tool"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  check         - Run all checks (fmt, vet, test)"
	@echo "  example       - Run basic example"
	@echo "  test-cli      - Test CLI with sample data"
	@echo "  mod-tidy      - Tidy go modules"
	@echo "  help          - Show this help"
