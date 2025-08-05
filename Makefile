# Page Insight Tool Makefile

.PHONY: build test run clean docker-build docker-run help

# Binary name and location
BINARY_NAME=page-insight-tool
BINARY_DIR=app/.bin

# Default target
all: build

# Build the application
build:
	@echo "Building Page Insight Tool..."
	@mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_DIR)/$(BINARY_NAME) app/cmd/page-insight-tool.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

test-app:
	./test_app.sh
# Run the application
run: build
	@echo "Starting Page Insight Tool..."
	./$(BINARY_DIR)/$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe
	rm -rf $(BINARY_DIR)
	rm -rf dist/
	rm -f *.log
	rm -rf tmp/

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t page-insight-tool .

# Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	docker run -p 8080:8080 page-insight-tool

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Show help
help:
	@echo "Page Insight Tool - Available commands:"
	@echo "  build        - Build the application"
	@echo "  test         - Run tests"
	@echo "  test-app     - test app running"
	@echo "  run          - Build and run the application"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Build and run Docker container"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  help         - Show this help"
