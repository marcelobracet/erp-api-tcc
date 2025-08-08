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
	# TODO: Add migration commands

seed:
	@echo "Seeding database..."
	# TODO: Add seed commands 