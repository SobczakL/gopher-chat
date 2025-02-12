# Variables
APP_NAME = gopher-chat
CMD_DIR = ./cmd/server

# Go commands
GO = go
BUILD_DIR = ./bin
MAIN_FILE = $(CMD_DIR)/main.go

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building the application..."
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "Build complete. Executable created at $(BUILD_DIR)/$(APP_NAME)."

# Run the application
.PHONY: run
run:
	@echo "Running the application..."
	$(BUILD_DIR)/$(APP_NAME)

# Test the application
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test ./... -v

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GO) mod tidy

# Help menu
.PHONY: help
help:
	@echo "Makefile for $(APP_NAME)"
	@echo "Usage:"
	@echo "  make build      Build the application"
	@echo "  make run        Run the application"
	@echo "  make test       Run tests"
	@echo "  make clean      Clean build artifacts"
	@echo "  make fmt        Format the code"
	@echo "  make deps       Install dependencies"
	@echo "  make help       Show this help menu"
