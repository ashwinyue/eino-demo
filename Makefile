SHELL := /bin/bash

SERVER_DIR := examples/go/orders-cs-adk/cmd/server
MCP_DIR    := examples/go/orders-cs-adk/mcpserver
CLI_DIR    := examples/go/orders-cs-adk
PKG_DIR    := examples/go/orders-cs-adk

.PHONY: help deps build test run-server run-mcp run-cli env

help:
	@echo "Available targets:"
	@echo "  deps        - tidy and fetch Go dependencies"
	@echo "  build       - build all packages"
	@echo "  test        - run unit/integration tests"
	@echo "  run-server  - start Gin HTTP server (:8080)"
	@echo "  run-mcp     - start local MCP-style mock (:8000)"
	@echo "  run-cli     - start interactive CLI (stdin)"
	@echo "  env         - copy example config to active config if missing"

deps:
	cd $(PKG_DIR) && go mod tidy

build:
	cd $(PKG_DIR) && go build ./...

test:
	cd $(PKG_DIR) && go test ./...

run-server:
	cd $(SERVER_DIR) && go run .

run-mcp:
	cd $(MCP_DIR) && go run .

run-cli:
	cd $(CLI_DIR) && go run .

env:
	@if [ ! -f $(PKG_DIR)/config.yaml ]; then \
		cp $(PKG_DIR)/config.example.yaml $(PKG_DIR)/config.yaml; \
		echo "Created $(PKG_DIR)/config.yaml from example"; \
	else \
		echo "$(PKG_DIR)/config.yaml already exists"; \
	fi

