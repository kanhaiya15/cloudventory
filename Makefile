.PHONY: help build run test lint clean docker-build docker-run migrate migrate-down deps install

# Variables
BINARY_NAME=aws-inventory-system
DOCKER_IMAGE=aws-inventory-system
DOCKER_TAG=latest
DATABASE_URL=postgres://postgres:password@localhost:5432/aws_inventory?sslmode=disable

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
deps: ## Download Go dependencies
	go mod download
	go mod tidy

build: ## Build the application
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/$(BINARY_NAME) ./cmd

run: ## Run the application locally
	go run ./cmd -db-url="$(DATABASE_URL)" -region=us-east-1 -parallel=true

install: ## Install the application
	go install ./cmd

# Testing
test: ## Run all tests
	go test -v ./...

test-integration: ## Run integration tests
	go test -v -tags=integration ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to coverage.html"

lint: ## Run linter
	golangci-lint run

# Database
migrate: ## Run database migrations up
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path migrations -database "$(DATABASE_URL)" up; \
	else \
		docker run --rm --network host -v $(PWD)/migrations:/migrations migrate/migrate:v4.16.2 \
			-path /migrations -database "$(DATABASE_URL)" up; \
	fi

migrate-down: ## Run database migrations down
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path migrations -database "$(DATABASE_URL)" down; \
	else \
		docker run --rm --network host -v $(PWD)/migrations:/migrations migrate/migrate:v4.16.2 \
			-path /migrations -database "$(DATABASE_URL)" down; \
	fi

migrate-create: ## Create a new migration file (usage: make migrate-create name=migration_name)
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a migration name. Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@if command -v migrate >/dev/null 2>&1; then \
		migrate create -ext sql -dir migrations -seq $(name); \
	else \
		docker run --rm -v $(PWD)/migrations:/migrations migrate/migrate:v4.16.2 \
			create -ext sql -dir /migrations -seq $(name); \
	fi

# Docker
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run application in Docker
	docker-compose up -d postgres
	sleep 10
	docker-compose run --rm migrate
	docker-compose up aws-inventory

docker-dev: ## Start development environment with Docker Compose
	docker-compose up -d postgres
	sleep 10
	docker-compose --profile migration run --rm migrate
	@echo "Database is ready. Run 'make run' to start the application."

docker-stop: ## Stop Docker containers
	docker-compose down

docker-clean: ## Clean up Docker containers and images
	docker-compose down -v --remove-orphans
	docker system prune -f

# Local development
dev-setup: ## Set up local development environment
	@echo "Setting up development environment..."
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "Error: Docker is required but not installed."; \
		exit 1; \
	fi
	@if ! command -v docker-compose >/dev/null 2>&1; then \
		echo "Error: Docker Compose is required but not installed."; \
		exit 1; \
	fi
	@echo "Starting PostgreSQL..."
	docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	sleep 15
	@echo "Running migrations..."
	make migrate
	@echo "Development environment ready!"
	@echo "You can now run 'make run' to start the application."

# Cleanup
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

# Production
deploy: ## Build and deploy (customize as needed)
	@echo "Building production image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Build complete. Image: $(DOCKER_IMAGE):$(DOCKER_TAG)"

# AWS credentials check
check-aws: ## Check AWS credentials
	@if [ -z "$$AWS_ACCESS_KEY_ID" ] && [ ! -f ~/.aws/credentials ]; then \
		echo "Warning: No AWS credentials found."; \
		echo "Please set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables"; \
		echo "or configure AWS credentials using 'aws configure'"; \
	else \
		echo "AWS credentials found."; \
	fi

# All-in-one setup
setup: dev-setup check-aws ## Complete setup for development environment

# Show current configuration
status: ## Show current configuration
	@echo "=== Configuration ==="
	@echo "Binary name: $(BINARY_NAME)"
	@echo "Docker image: $(DOCKER_IMAGE):$(DOCKER_TAG)"
	@echo "Database URL: $(DATABASE_URL)"
	@echo ""
	@echo "=== Services Status ==="
	@docker-compose ps 2>/dev/null || echo "Docker Compose not running"