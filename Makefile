.PHONY: build test clean lint vet help run

BINARY_NAME=moneyback
BINARY_PATH=./bin/$(BINARY_NAME)

help:
	@echo "money-backword Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  make build       - Build the CLI binary"
	@echo "  make test        - Run all tests"
	@echo "  make test-v      - Run tests with verbose output"
	@echo "  make test-cov    - Run tests with coverage report"
	@echo "  make vet         - Run go vet"
	@echo "  make lint        - Run golangci-lint"
	@echo "  make run         - Build and run the CLI"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make deps        - Download dependencies"

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	@go build -o $(BINARY_PATH) ./cmd/moneyback
	@echo "Built: $(BINARY_PATH)"

run: build
	@$(BINARY_PATH) -help

test:
	@go test -race ./...

test-v:
	@go test -v -race ./...

test-cov:
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

vet:
	@go vet ./...

lint:
	@golangci-lint run ./... 2>/dev/null || echo "golangci-lint not installed, skipping"

clean:
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Cleaned"

deps:
	@go mod download
	@go mod verify

fmt:
	@go fmt ./...

fmt-check:
	@if [ -n "$$(go fmt ./...)" ]; then echo "Code needs formatting"; exit 1; fi

.DEFAULT_GOAL := help
