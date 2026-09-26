# ─────────────────────────────────────────────────────────────────────────────
#  Makefile — build, test and documentation targets
#
#  The binary name and main package come from docsmith.toml, so renaming the
#  CLI only means editing that file.
# ─────────────────────────────────────────────────────────────────────────────

BINARY  := $(shell awk -F'"' '/^binary *=/ {print $$2; exit}' docsmith.toml)
MAIN    := $(shell awk -F'"' '/^main *=/ {print $$2; exit}' docsmith.toml)
PKGS    := ./...
GOBIN   := $(shell go env GOPATH)/bin

# freeze and golangci-lint are usually installed into GOBIN.
export PATH := $(GOBIN):$(PATH)

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -X main.version=$(VERSION)

LINT := $(shell command -v golangci-lint 2>/dev/null || echo "go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest")

# ANSI colour codes (expanded by Make, interpreted by printf)
GREEN   := \033[32m
YELLOW  := \033[33m
RED     := \033[31m
CYAN    := \033[36m
BOLD    := \033[1m
DIM     := \033[2m
RESET   := \033[0m

# ─────────────────────────────────────────────────────────────────────────────

.DEFAULT_GOAL := help

.PHONY: help build test lint fmt docs banner diagrams examples examples-check generate docs-all tools

help: ## Show available targets
	@printf "\n$(BOLD)$(BINARY)$(RESET) — development commands\n\n"
	@grep -E '^[a-zA-Z/_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	  | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-15s$(RESET) %s\n", $$1, $$2}'
	@printf "\n"

# ── Build / test / lint ──────────────────────────────────────────────────────

build: ## Build the CLI into ./bin/
	@printf "\n$(BOLD)$(CYAN)── Build$(RESET)\n"
	@go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) $(MAIN) \
	  && printf "  $(GREEN)✓$(RESET)  $(BOLD)$(BINARY)$(RESET)  $(DIM)→ bin/$(BINARY)$(RESET)\n"

test: ## Run the test suite, including the example scripts
	@printf "\n$(BOLD)$(CYAN)── Test$(RESET)\n"
	@go test -count=1 $(PKGS)

lint: ## Run golangci-lint
	@printf "\n$(BOLD)$(CYAN)── Lint$(RESET)\n"
	@$(LINT) run $(PKGS) && printf "  $(GREEN)✓$(RESET)  $(BOLD)golangci-lint$(RESET)\n"

fmt: ## Format all Go source files with gofmt
	@gofmt -l -w .

# ── Documentation ────────────────────────────────────────────────────────────

docs: ## Regenerate the CLI reference and the project layout
	@printf "\n$(BOLD)$(CYAN)── Docs$(RESET)\n"
	@go run ./internal/tools/gendocs

banner: ## Regenerate the banner SVGs from docsmith.toml
	@printf "\n$(BOLD)$(CYAN)── Banner$(RESET)\n"
	@go run ./internal/tools/genbanner

diagrams: ## Regenerate diagram SVGs from docs/diagrams/*.toml
	@printf "\n$(BOLD)$(CYAN)── Diagrams$(RESET)\n"
	@go run ./internal/tools/gendiagram

examples: ## Run the example scripts, then render them into the docs
	@printf "\n$(BOLD)$(CYAN)── Examples$(RESET)\n"
	@go test -count=1 -run TestScript $(MAIN) \
	  && printf "  $(GREEN)✓$(RESET)  $(BOLD)testscript$(RESET)  $(DIM)all examples passed$(RESET)\n"
	@if ! command -v freeze >/dev/null 2>&1; then \
	   printf "  $(YELLOW)⚠$(RESET)  $(BOLD)freeze$(RESET)  $(DIM)not installed — code blocks only; run: make tools$(RESET)\n"; \
	 fi
	@go run ./internal/tools/genexamples

examples-check: ## Fail when the rendered examples in the docs are out of date (CI)
	@printf "\n$(BOLD)$(CYAN)── Examples (check)$(RESET)\n"
	@go test -count=1 -run TestScript $(MAIN) \
	  && printf "  $(GREEN)✓$(RESET)  $(BOLD)testscript$(RESET)  $(DIM)all examples passed$(RESET)\n"
	@go run ./internal/tools/genexamples
	@git diff --exit-code -- '*.md' \
	  && printf "  $(GREEN)✓$(RESET)  $(BOLD)drift check$(RESET)  $(DIM)docs up to date$(RESET)\n" \
	  || (printf "  $(RED)✗$(RESET)  $(BOLD)drift check$(RESET)  $(DIM)run: make examples$(RESET)\n"; exit 1)

generate: banner diagrams docs examples ## Regenerate every generated file

docs-all: generate ## Alias for generate

tools: ## Install freeze and golangci-lint into GOBIN
	@go install github.com/charmbracelet/freeze@latest
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
