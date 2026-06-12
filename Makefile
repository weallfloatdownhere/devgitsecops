.PHONY: build install clean test run help

# Variables
BINARY_NAME=devgitsecops
MAIN_PATH=main.go
BUILD_DIR=bin
GO=go
GOFLAGS=-v

# Build information
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# LDFLAGS
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

help: ## Display this help screen
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-all: ## Build for all platforms
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "Multi-platform build complete"

install: build ## Build and install to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(LDFLAGS) $(MAIN_PATH)
	@echo "Installation complete"

clean: ## Remove build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@$(GO) clean
	@echo "Clean complete"

test: ## Run tests
	@echo "Running tests..."
	$(GO) test -v ./...

run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "Dependencies updated"

fmt: ## Format code
	@echo "Formatting code..."
	$(GO) fmt ./...

lint: ## Run linter
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GO) vet ./...

check: fmt vet lint test ## Run all checks (fmt, vet, lint, test)
