set shell := ["bash", "-c"]

# Default target when running `just` without arguments
default: help

# Run tests then build
all: test build

# Build the binary
build:
    go build -o bin/mux-session -v .

# Build and run the binary
run: build
    ./bin/mux-session

# Run tests
test:
    go test -v ./...

# Run end-to-end tests
e2e:
    go test -count 1 -v ./e2e/...

# Clean build artifacts
clean:
    go clean
    rm -f bin/mux-session
    rm -f bin/mux-session_unix

# Download and tidy dependencies
deps:
    go mod download
    go mod tidy

# Cross compilation for linux
build-linux:
    env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/mux-session_unix -v .

# Run directly without building binary
dev:
    go run main.go

# Install binary globally with go install
install:
    go install .

# Format the codebase using go fmt
fmt:
    go fmt ./...

# Display available commands
help:
    @just --list
