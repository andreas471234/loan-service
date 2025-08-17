# Makefile for Loan Service

# Variables
BINARY_NAME=loan-service
BUILD_DIR=build
MAIN_PATH=cmd/server/main.go

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(shell git describe --tags --always --dirty)"

.PHONY: all build clean test coverage run docker-build docker-run help

# Default target
all: clean build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	$(GOTEST) -v ./internal/...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v ./tests/integration/...

# Run all tests
test:
	@echo "Running all tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=tools/coverage.out ./...
	$(GOCMD) tool cover -html=tools/coverage.out -o tools/coverage.html
	@echo "Coverage report generated: tools/coverage.html"

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run $(MAIN_PATH)

# Run the application with race detection
run-race:
	@echo "Running $(BINARY_NAME) with race detection..."
	$(GOCMD) run -race $(MAIN_PATH)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	$(GOCMD) vet ./...

# Run all checks (format, lint, test)
check: fmt lint test

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) -f deployments/Dockerfile .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(BINARY_NAME)

# Docker compose up
docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose -f deployments/docker-compose.yml up -d

# Docker compose down
docker-down:
	@echo "Stopping services with Docker Compose..."
	docker-compose -f deployments/docker-compose.yml down

# Generate API documentation (if using swagger)
docs:
	@echo "Generating API documentation..."
	# Add swagger generation commands here if needed

# Database migration
migrate:
	@echo "Running database migrations..."
	$(GOCMD) run $(MAIN_PATH) migrate

# Show help
help:
	@echo "Available targets:"
	@echo "  build       - Build the application"
	@echo "  clean       - Clean build artifacts"
	@echo "  test        - Run all tests"
	@echo "  test-unit   - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  coverage    - Run tests with coverage report"
	@echo "  run         - Run the application"
	@echo "  run-race    - Run with race detection"
	@echo "  deps        - Install dependencies"
	@echo "  fmt         - Format code"
	@echo "  lint        - Run linter"
	@echo "  check       - Run format, lint, and test"
	@echo "  docker-build- Build Docker image"
	@echo "  docker-run  - Run Docker container"
	@echo "  docker-up   - Start with Docker Compose"
	@echo "  docker-down - Stop Docker Compose services"
	@echo "  docs        - Generate API documentation"
	@echo "  migrate     - Run database migrations"
	@echo "  help        - Show this help message"