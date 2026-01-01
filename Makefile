# Makefile for Beto Go Application

# Variables
APP_NAME=beto
BINARY_NAME=beto
VERSION=1.0.0
BUILD_DIR=build
DOCKER_IMAGE=beto:latest
GO_VERSION=1.25

# Default target
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
.PHONY: build
build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) .

.PHONY: build-windows
build-windows: ## Build for Windows
	@echo "Building $(BINARY_NAME) for Windows..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME).exe .

.PHONY: build-mac
build-mac: ## Build for macOS
	@echo "Building $(BINARY_NAME) for macOS..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-mac .

.PHONY: build-all
build-all: build build-windows build-mac ## Build for all platforms

# Development targets
.PHONY: run
run: ## Run the application
	@echo "Running $(APP_NAME)..."
	@go run .

.PHONY: dev
dev: ## Run with live reload (requires air)
	@echo "Starting development server with live reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Installing..."; \
		go install github.com/air-verse/air@latest; \
		air; \
	fi

.PHONY: watch
watch: dev ## Alias for dev

# Testing targets
.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	@go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-race
test-race: ## Run tests with race condition detection
	@echo "Running tests with race detection..."
	@go test -race -v ./...

.PHONY: bench
bench: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Code quality targets
.PHONY: lint
lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@gofumpt -w .

.PHONY: imports
imports: ## Fix imports
	@echo "Organizing imports..."
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
	else \
		echo "goimports not installed. Installing..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
		goimports -w .; \
	fi

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

.PHONY: staticcheck
staticcheck: ## Run staticcheck
	@echo "Running staticcheck..."
	@if command -v staticcheck > /dev/null; then \
		staticcheck ./...; \
	else \
		echo "staticcheck not installed. Installing..."; \
		go install honnef.co/go/tools/cmd/staticcheck@latest; \
		staticcheck ./...; \
	fi

.PHONY: check
check: fmt imports vet lint staticcheck ## Run all code quality checks

# Dependency management
.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

.PHONY: deps-verify
deps-verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	@go mod verify

.PHONY: tidy
tidy: ## Tidy up dependencies
	@echo "Tidying up dependencies..."
	@go mod tidy

# Database targets (if using database)
.PHONY: db-migrate
db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	@# Add your migration command here
	@echo "No migrations configured"

.PHONY: db-rollback
db-rollback: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	@# Add your rollback command here
	@echo "No rollback configured"

# Docker targets
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run --rm -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

.PHONY: docker-push
docker-push: docker-build ## Build and push Docker image
	@echo "Pushing Docker image..."
	@docker push $(DOCKER_IMAGE)

# Installation targets
.PHONY: install
install: ## Install the binary
	@echo "Installing $(BINARY_NAME)..."
	@go install -ldflags="-w -s -X main.version=$(VERSION)" .

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install mvdan.cc/gofumpt@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest

# Generate targets
.PHONY: generate
generate: ## Run go generate
	@echo "Running go generate..."
	@go generate ./...

.PHONY: mocks
mocks: ## Generate mocks
	@echo "Generating mocks..."
	@if command -v mockgen > /dev/null; then \
		echo "Generating mocks with mockgen..."; \
	else \
		echo "mockgen not installed. Installing..."; \
		go install github.com/golang/mock/mockgen@latest; \
	fi

# Documentation targets
.PHONY: docs
docs: ## Generate documentation
	@echo "Generating documentation..."
	@go doc -all ./... > docs/api.md

.PHONY: serve-docs
serve-docs: ## Serve documentation
	@echo "Serving documentation on :6060..."
	@if command -v godoc > /dev/null; then \
		godoc -http=:6060; \
	else \
		echo "godoc not installed. Installing..."; \
		go install golang.org/x/tools/cmd/godoc@latest; \
		godoc -http=:6060; \
	fi

# Profiling targets
.PHONY: profile-cpu
profile-cpu: build ## Profile CPU usage
	@echo "Running CPU profiling..."
	@./$(BUILD_DIR)/$(BINARY_NAME) -cpuprofile=cpu.prof &
	@sleep 10
	@killall $(BINARY_NAME) || true
	@go tool pprof cpu.prof

.PHONY: profile-mem
profile-mem: build ## Profile memory usage
	@echo "Running memory profiling..."
	@./$(BUILD_DIR)/$(BINARY_NAME) -memprofile=mem.prof &
	@sleep 10
	@killall $(BINARY_NAME) || true
	@go tool pprof mem.prof

# Security targets
.PHONY: security
security: ## Run security checks
	@echo "Running security checks..."
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Installing..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi

# Clean targets
.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@rm -f cpu.prof mem.prof
	@go clean

.PHONY: clean-deps
clean-deps: ## Clean dependency cache
	@echo "Cleaning dependency cache..."
	@go clean -modcache

# Release targets
.PHONY: tag
tag: ## Create a git tag
	@echo "Creating tag v$(VERSION)..."
	@git tag -a v$(VERSION) -m "Release version $(VERSION)"
	@git push origin v$(VERSION)

.PHONY: release
release: clean check test build-all ## Prepare release
	@echo "Release $(VERSION) ready in $(BUILD_DIR)/"

# Development workflow
.PHONY: setup
setup: deps install-tools ## Setup development environment
	@echo "Development environment setup complete!"
	@echo "Run 'make dev' to start development server"

.PHONY: ci
ci: check test-coverage test-race ## Run CI pipeline locally

# Environment setup
.PHONY: env
env: ## Copy .env.example to .env
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env file created from .env.example"; \
	else \
		echo ".env file already exists"; \
	fi

# Status and info
.PHONY: status
status: ## Show project status
	@echo "Project: $(APP_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Go Version: $(shell go version)"
	@echo "Dependencies:"
	@go list -m all | head -10

.PHONY: version
version: ## Show version
	@echo $(VERSION)

# Default make target
.DEFAULT_GOAL := help
