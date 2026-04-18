.PHONY: dev build test lint run

# Development
dev:
	@go run cmd/server/main.go

# Build
build:
	@go build -ldflags="-s -w" -o bin/server cmd/server/main.go

# Run production
run: build
	@./bin/server

# Test
test:
	@go test -v -race ./...

# Test with coverage
test-cover:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Lint
lint:
	@golangci-lint run ./...

# Tidy
tidy:
	@go mod tidy

# Generate
generate:
	@go generate ./...
