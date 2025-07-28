# Build configuration
BINARY_SERVER=bin/nanorediz-server
BINARY_CLIENT=bin/nanorediz-client
BINARY_WEB=bin/nanorediz-web

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.version=$(shell git describe --tags --always --dirty) -X main.buildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"

.PHONY: all build clean test deps help server client web

## Build all binaries
all: deps build

## Build all binaries
build: server client web

## Build server binary
server:
	@echo "Building server..."
	@mkdir -p bin
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_SERVER) ./cmd/server

## Build client binary
client:
	@echo "Building client..."
	@mkdir -p bin
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_CLIENT) ./cmd/client

## Build web binary
web:
	@echo "Building web interface..."
	@mkdir -p bin
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_WEB) ./web

## Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/
	rm -f coverage.out coverage.html

## Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"; \
	fi

## Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

## Start a development cluster
dev-cluster: server
	@echo "Starting development cluster..."
	@./bin/nanorediz-server -host 127.0.0.1 -port 8080 &
	@sleep 2
	@./bin/nanorediz-server -host 127.0.0.1 -port 8081 -contact-host 127.0.0.1 -contact-port 8080 &
	@sleep 2
	@./bin/nanorediz-server -host 127.0.0.1 -port 8082 -contact-host 127.0.0.1 -contact-port 8080 &
	@echo "Cluster started on ports 8080, 8081, 8082"

## Stop development cluster
stop-cluster:
	@echo "Stopping development cluster..."
	@pkill -f "nanorediz-server" || true

## Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GOCMD) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

## Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	@if command -v protoc >/dev/null 2>&1; then \
		protoc --go_out=. --go_opt=paths=source_relative \
		       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		       protos/*.proto; \
	else \
		echo "protoc not installed. Please install Protocol Buffers compiler"; \
	fi

## Show help
help:
	@echo "Available commands:"
	@grep -E '^## .*' Makefile | sed 's/## /  /'
	@echo ""
	@echo "Environment variables:"
	@echo "  NANOREDIZ_HOST              Server host (default: 0.0.0.0)"
	@echo "  NANOREDIZ_PORT              Server port (default: 8080)"
	@echo "  NANOREDIZ_LOG_LEVEL         Log level (debug, info, warn, error)"
	@echo "  NANOREDIZ_GRPC_TIMEOUT      gRPC timeout (default: 30s)"
	@echo "  NANOREDIZ_SHUTDOWN_TIMEOUT  Graceful shutdown timeout (default: 10s)"