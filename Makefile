.PHONY: build run test clean deps lint help

# Default target
help:
	@echo "Available commands:"
	@echo "  build    - Build the application"
	@echo "  run      - Run the application"
	@echo "  test     - Run tests"
	@echo "  clean    - Clean build artifacts"
	@echo "  deps     - Install dependencies"
	@echo "  lint     - Run linter"
	@echo "  api-test - Run API tests"

# Build the application
build:
	go build -o loan-service .

# Run the application
run:
	go run .

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f loan-service
	rm -f loan_service.db

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run linter
lint:
	golangci-lint run

# Run API tests
api-test:
	./test_api.sh

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest 