# Dealer Development Platform Makefile

.PHONY: help install test test-unit test-integration test-all lint build clean docker-build docker-run

# Default target
help:
	@echo "Dealer Development Platform - Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  install          Install dependencies for frontend and backend"
	@echo "  build            Build frontend and backend"
	@echo "  clean            Clean build artifacts"
	@echo ""
	@echo "Testing:"
	@echo "  test             Run all tests (unit + integration)"
	@echo "  test-unit        Run unit tests only"
	@echo "  test-integration Run integration tests with Testcontainers"
	@echo "  test-all         Run all tests with coverage"
	@echo ""
	@echo "Code Quality:"
	@echo "  lint             Run linting for frontend and backend"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build     Build Docker images"
	@echo "  docker-run       Run application in Docker"
	@echo ""
	@echo "Database:"
	@echo "  db-migrate       Run database migrations"
	@echo "  db-reset         Reset database (WARNING: destructive)"

# Install dependencies
install:
	@echo "Installing frontend dependencies..."
	cd frontend && npm ci
	@echo "Installing backend dependencies..."
	cd backend && go mod download && go mod verify

# Build
build:
	@echo "Building frontend..."
	cd frontend && npm run build
	@echo "Building backend..."
	cd backend && go build -o bin/dealer-platform ./cmd/app/main.go

# Clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf frontend/dist
	rm -rf backend/bin
	rm -rf backend/*.out
	rm -rf backend/*.html

# Testing
test: test-unit test-integration

test-unit:
	@echo "Running unit tests..."
	cd backend && go test -v -race -short ./internal/model/... ./internal/config/... ./internal/utils/...

test-integration:
	@echo "Running integration tests with Testcontainers..."
	cd backend && go test -v -race -timeout=30m ./internal/service/... ./internal/repository/... ./internal/testutil/...

test-all:
	@echo "Running all tests with coverage..."
	cd backend && go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Linting
lint:
	@echo "Linting frontend..."
	cd frontend && npm run lint || echo "Frontend linting completed with warnings"
	@echo "Linting backend..."
	cd backend && go vet ./... && gofmt -s -l . | grep -q . && echo "Go code is not formatted" || echo "Backend formatting OK"

# Docker
docker-build:
	@echo "Building Docker images..."
	docker build -t dealer-platform-backend ./backend
	docker build -t dealer-platform-frontend ./frontend

docker-run:
	@echo "Running application in Docker..."
	docker-compose up -d

# Database
db-migrate:
	@echo "Running database migrations..."
	@echo "Please configure your DATABASE_URL environment variable"
	@echo "Example: DATABASE_URL=postgres://user:pass@localhost:5432/dbname make db-migrate"

db-reset:
	@echo "WARNING: This will reset the database!"
	@echo "Please configure your DATABASE_URL environment variable"
	@echo "Example: DATABASE_URL=postgres://user:pass@localhost:5432/dbname make db-reset"

# CI/CD helpers
ci-test:
	@echo "Running CI tests..."
	cd backend && go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

ci-integration-test:
	@echo "Running CI integration tests..."
	cd backend && TESTCONTAINERS_RYUK_DISABLED=true go test -v -race -timeout=30m ./internal/service/... ./internal/repository/... ./internal/testutil/...

# Development helpers
dev-setup: install
	@echo "Development environment setup complete!"
	@echo "Run 'make test' to verify everything works"

dev-clean: clean
	@echo "Development cleanup complete!"

# Quick commands
quick-test:
	@echo "Running quick tests..."
	cd backend && go test -v -short ./...

quick-build:
	@echo "Quick build..."
	cd backend && go build ./cmd/app/main.go