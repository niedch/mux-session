# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=mux-session
BINARY_UNIX=$(BINARY_NAME)_unix

.PHONY: all build clean test run deps help install e2e

all: test build

build: 
	$(GOBUILD) -o bin/$(BINARY_NAME) -v .

test:
	$(GOTEST) -v ./...

e2e:
	$(GOTEST) -count 1 -v ./e2e/...

run:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v .
	./bin/$(BINARY_NAME)

clean: 
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)
	rm -f bin/$(BINARY_UNIX)

deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_UNIX) -v .

# Development
dev:
	$(GOCMD) run main.go

install:
	$(GOCMD) install .

help:
	@echo "Available commands:"
	@echo "  make build    - Build the binary"
	@echo "  make run      - Build and run the binary"
	@echo "  make test     - Run tests"
	@echo "  make e2e      - Run end-to-end tests"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make deps     - Download and tidy dependencies"
	@echo "  make dev      - Run directly without building binary"
	@echo "  make install  - Install binary globally with go install"
	@echo "  make all      - Run tests then build"
