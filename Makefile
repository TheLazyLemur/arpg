# ARPG Game Makefile

# Variables
BINARY_NAME=arpg
BUILD_DIR=build
CMD_DIR=cmd/arpg
MAIN_FILE=$(CMD_DIR)/main.go

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Build flags
LDFLAGS=-ldflags "-s -w"
BUILD_FLAGS=-v

# Default target
.PHONY: all
all: clean build

# Build the project
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the game
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Tidy dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

# Test the project
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Build for different platforms
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)

.PHONY: build-mac
build-mac:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)

# Build for all platforms
.PHONY: build-all
build-all: build-linux build-windows build-mac

# Development setup
.PHONY: setup
setup: deps tidy
	@echo "Development setup complete"

# Release build
.PHONY: release
release: clean fmt test build-all
	@echo "Release build complete"

# Development run with hot reload (requires air)
.PHONY: dev
dev:
	@if command -v air > /dev/null; then \
		echo "Starting development server with hot reload..."; \
		air; \
	else \
		echo "Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Running without hot reload..."; \
		$(MAKE) run; \
	fi

# Install air for hot reload
.PHONY: install-air
install-air:
	@echo "Installing air for hot reload..."
	$(GOGET) github.com/cosmtrek/air@latest

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build       - Build the project"
	@echo "  run         - Build and run the game"
	@echo "  clean       - Clean build artifacts"
	@echo "  fmt         - Format code"
	@echo "  tidy        - Tidy dependencies"
	@echo "  deps        - Download dependencies"
	@echo "  test        - Run tests"
	@echo "  build-linux - Build for Linux"
	@echo "  build-windows - Build for Windows"
	@echo "  build-mac   - Build for macOS"
	@echo "  build-all   - Build for all platforms"
	@echo "  setup       - Set up development environment"
	@echo "  release     - Create release build"
	@echo "  dev         - Run with hot reload (requires air)"
	@echo "  install-air - Install air for hot reload"
	@echo "  help        - Show this help message"