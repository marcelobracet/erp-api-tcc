.PHONY: help build test clean run docker-build docker-run k8s-deploy

help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make test         - Run all tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make run          - Run the application"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run with Docker Compose"
	@echo "  make k8s-deploy   - Deploy to Kubernetes"

build:
	@echo "Building application..."
	go build -o bin/erp-api cmd/api/main.go

test:
	@echo "Running tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

run:
	@echo "Running application..."
	go run cmd/api/main.go

docker-build:
	@echo "Building Docker image..."
	docker build -t erp-api:latest .

docker-run:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

docker-stop:
	@echo "Stopping Docker Compose services..."
	docker-compose down

k8s-deploy:
	@echo "Deploying to Kubernetes..."
	kubectl apply -f k8s/

deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

mocks:
	@echo "Generating mocks..."
	mockery --all --output=./mocks

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Linting code..."
	golangci-lint run

migrate:
	@echo "Running database migrations..."
	@docker-compose exec postgres psql -U erp_user -d erp_db -f /docker-entrypoint-initdb.d/001_create_users_table.sql

seed:
	@echo "Seeding database..."
	@echo "Seed commands to be implemented"

db-setup:
	@echo "Setting up database..."
	@docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 10
	@echo "Database setup complete"

dev-setup: deps build db-setup
	@echo "Development environment setup complete"

dev: dev-setup
	@echo "Starting development environment..."
	@docker-compose up

dev-stop:
	@echo "Stopping development environment..."
	@docker-compose down

logs:
	@echo "Showing logs..."
	@docker-compose logs -f

db-connect:
	@echo "Connecting to PostgreSQL..."
	@docker-compose exec postgres psql -U erp_user -d erp_db

pgadmin:
	@echo "PgAdmin available at http://localhost:5050"
	@echo "Email: admin@erp.com"
	@echo "Password: admin123"

# Security and validation
security-scan:
	@echo "Running security scan..."
	@gosec ./...

validate: lint security-scan test
	@echo "All validations passed!"

# CI/CD commands
ci-test: deps test security-scan
	@echo "CI tests completed"

ci-build: deps build
	@echo "CI build completed"

# Deployment commands
deploy-qa:
	@echo "Deploying to QA environment..."
	@chmod +x scripts/deploy-qa.sh
	@./scripts/deploy-qa.sh

deploy-prod:
	@echo "Deploying to Production environment..."
	@chmod +x scripts/deploy-prod.sh
	@./scripts/deploy-prod.sh

# Release management
create-release:
	@echo "Creating new release..."
	@git tag -a v$(shell date +%Y%m%d.%H%M) -m "Release $(shell date +%Y-%m-%d)"
	@git push --tags

# Kubernetes commands
k8s-apply-qa:
	@echo "Applying QA Kubernetes configurations..."
	@kubectl apply -f k8s/qa/ -n qa

k8s-apply-prod:
	@echo "Applying Production Kubernetes configurations..."
	@kubectl apply -f k8s/prod/ -n production

k8s-status-qa:
	@echo "QA Environment Status:"
	@kubectl get pods,svc,ing -n qa

k8s-status-prod:
	@echo "Production Environment Status:"
	@kubectl get pods,svc,ing -n production 