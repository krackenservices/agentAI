# Makefile for building and testing the application
APPNAME=agentAI
BINARY_API=bin/$(APPNAME)-api
BINARY_CLI=bin/$(APPNAME)-cli

TOOLS_DIR := pkg/tools
TOOL_NAMES := $(notdir $(wildcard $(TOOLS_DIR)/*))
BINARY_TOOLS := $(patsubst %,bin/tools/$(APPNAME)-%, $(TOOL_NAMES))

.PHONY: all build clean docker build-tools api cli tools test

all: build build-tools

build: api cli

debug-api: clean-api api
	dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec $(BINARY_API)

api:
	@echo "Building API server..."
	go build -o $(BINARY_API) ./cmd/api

swagger:
	swag init -d cmd/api,internal/handlers,internal/config,internal/routes,internal/server,internal/toolmodel,internal/toolregistry -o swdocs

cli:
	@echo "Building CLI tool..."
	go build -o $(BINARY_CLI) ./cmd/cli

# Pattern rule to build each tool found in pkg/tools.
bin/tools/$(APPNAME)-%: $(TOOLS_DIR)/%
	@echo "Building tool $*..."
	mkdir -p $(dir $@)
	go build -o $@ ./$<

# Build all tools.
tools: $(BINARY_TOOLS)

copy-config:
	@echo "Copying configuration file to API binary directory..."
	cp config.yml.sample $(dir $(BINARY_API))/config.yml;

test:
	@echo "Running tests..."
	go test -v ./...

# Coverage target: run tests with coverage and output total coverage stats.
coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	@echo "Coverage summary:"
	@go tool cover -func=coverage.out

clean-api:
	@echo "Cleaning up API server..."
	rm -f $(BINARY_API)

clean-tools:
	@echo "Cleaning up API server..."
	rm -f $(BINARY_TOOLS)

clean:
	@echo "Cleaning up..."
	rm -rf bin

docker:
	@echo "Building Docker image..."
	docker build -t agentAI .
