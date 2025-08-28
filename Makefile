# Makefile
.PHONY: build run test clean lint help mongo-start mongo-stop

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=server
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PATH=cmd/server/main.go

# Build the application
build:
	$(GOBUILD) -o $(BINARY_PATH) -v $(MAIN_PATH)

# Run the application
run:
	$(GOCMD) run $(MAIN_PATH)

# Run tests
test:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

# Test coverage
coverage: test
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)
	rm -f coverage.out coverage.html

# Start MongoDB (macOS with Homebrew)
mongo-start:
	@echo "Starting MongoDB..."
	@if command -v brew >/dev/null 2>&1; then \
		brew services start mongodb/brew/mongodb-community || brew services start mongodb-community; \
	elif command -v systemctl >/dev/null 2>&1; then \
		sudo systemctl start mongod; \
	else \
		echo "Please start MongoDB manually"; \
	fi

# Stop MongoDB
mongo-stop:
	@echo "Stopping MongoDB..."
	@if command -v brew >/dev/null 2>&1; then \
		brew services stop mongodb/brew/mongodb-community || brew services stop mongodb-community; \
	elif command -v systemctl >/dev/null 2>&1; then \
		sudo systemctl stop mongod; \
	else \
		echo "Please stop MongoDB manually"; \
	fi

# Check MongoDB status
mongo-status:
	@echo "MongoDB Status:"
	@if command -v brew >/dev/null 2>&1; then \
		brew services list | grep mongodb; \
	elif command -v systemctl >/dev/null 2>&1; then \
		sudo systemctl status mongod --no-pager -l; \
	else \
		echo "Please check MongoDB status manually"; \
	fi

# Connect to MongoDB shell
mongo-shell:
	@echo "Connecting to MongoDB shell..."
	@mongosh mongodb://localhost:27017/urlshortener

# Lint code (install golangci-lint first: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Build for production
build-prod:
	CGO_ENABLED=0 $(GOBUILD) -a -ldflags '-w -s' -o $(BINARY_PATH) $(MAIN_PATH)

# Run development server with hot reload (install air first: go install github.com/cosmtrek/air@latest)
dev:
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "Air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Running without hot reload:"; \
		$(GOCMD) run $(MAIN_PATH); \
	fi

# Setup development environment
setup:
	@echo "Setting up development environment..."
	@cp .env.example .env
	@$(GOMOD) tidy
	@echo "Environment file created (.env)"
	@echo "Install MongoDB if not already installed:"
	@echo "  macOS: brew install mongodb-community"
	@echo "  Ubuntu: Follow MongoDB installation guide"
	@echo "Then run: make mongo-start"

# Help
help:
	@echo "Available commands:"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application"
	@echo "  test        - Run tests"
	@echo "  coverage    - Run tests with coverage"
	@echo "  clean       - Clean build artifacts"
	@echo "  mongo-start - Start MongoDB service"
	@echo "  mongo-stop  - Stop MongoDB service"
	@echo "  mongo-status- Check MongoDB status"
	@echo "  mongo-shell - Connect to MongoDB shell"
	@echo "  lint        - Lint code"
	@echo "  deps        - Download dependencies"
	@echo "  build-prod  - Build for production"
	@echo "  dev         - Run with hot reload (requires air)"
	@echo "  setup       - Setup development environment"