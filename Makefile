# Makefile for baseline-init

# Build variables
BINARY_NAME=baseline-init
VERSION?=dev
GIT_COMMIT=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X github.com/aguamala/baseline-init/cmd.Version=$(VERSION) -X github.com/aguamala/baseline-init/cmd.GitCommit=$(GIT_COMMIT) -X github.com/aguamala/baseline-init/cmd.BuildDate=$(BUILD_DATE)"

.PHONY: all build clean test test-coverage lint install help

all: build ## Build the project

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o $(BINARY_NAME) .

clean: ## Remove build artifacts
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@go clean

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linters
	@echo "Running linters..."
	@go fmt ./...
	@go vet ./...

install: ## Install the binary
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) .

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

run-check: build ## Build and run check command
	@./$(BINARY_NAME) check

run-setup: build ## Build and run setup command
	@./$(BINARY_NAME) setup --auto

help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
