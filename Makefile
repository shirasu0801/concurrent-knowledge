.PHONY: help build run test clean init-db

help: ## Show this help message
	@echo "Concurrent Knowledge - Makefile Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building Concurrent Knowledge..."
	go build -o bin/server.exe cmd/server/main.go
	@echo "Build complete: bin/server.exe"

run: ## Run the application
	@echo "Starting Concurrent Knowledge server..."
	go run cmd/server/main.go

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf data/*.db*
	@echo "Clean complete"

init-db: ## Initialize database with migrations
	@echo "Initializing database..."
	cd scripts && ./init_db.bat
	@echo "Database initialized"

deps: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed"

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...
	@echo "Vet complete"

dev: init-db run ## Initialize DB and run in development mode
