# Makefile for tablo CLI
SHELL := /bin/bash

# Derive version: if latest git tag (v*) exists use it; else use dev-<short-hash>; final fallback 'dev'. Override: make VERSION=... build
VERSION ?= $(shell (git describe --tags --abbrev=0 2>/dev/null || (h=$$(git rev-parse --short HEAD 2>/dev/null) && echo dev-$$h) || echo dev))
TAG_NAME := $(if $(filter v%,$(VERSION)),$(VERSION),v$(VERSION))
MIN_COVER ?= 70.0
OS_ARCHES = linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64
BINARY = tablo
BIN_DIR = bin
DIST_DIR = dist
PKG = ./cmd/tablo

# Default target
.PHONY: help
help: ## Show this help
	@echo "tablo Makefile targets"; echo; awk 'BEGIN {FS = ":.*##"; printf "Usage: make <target>\n\nTargets:\n"} /^[a-zA-Z0-9_.-]+:.*##/ { printf "  %-20s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the CLI binary (local OS/ARCH)
	@echo "Building $(BINARY) (version $(VERSION))"
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o $(BIN_DIR)/$(BINARY) $(PKG)
	@echo "Built: $(BIN_DIR)/$(BINARY)"

.PHONY: install
install: ## Install the binary into GOPATH/bin (or GOBIN)
	@echo "Installing $(BINARY) (version $(VERSION))"
	CGO_ENABLED=0 go install -trimpath -ldflags="-s -w -X main.version=$(VERSION)" $(PKG)

.PHONY: run
run: build ## Run the CLI with --help (example)
	./$(BIN_DIR)/$(BINARY) --help

.PHONY: tidy
tidy: ## Run go mod tidy (and verify no drift)
	go mod tidy
	@git diff --quiet -- go.mod go.sum || (echo "go.mod/go.sum changed; commit the updates" && exit 1)

.PHONY: lint
lint: ## Run linters (golangci-lint if available, else go vet/gofmt/staticcheck)
	@if command -v golangci-lint >/dev/null 2>&1; then \
	  echo "Running golangci-lint"; \
	  golangci-lint run ./...; \
	else \
	  echo "golangci-lint not found; running go vet"; \
	  go vet ./...; \
	fi
	@echo "Checking formatting (gofmt)"; \
	fmt_out=$$(gofmt -l . | grep -v '^vendor/' || true); \
	if [ -n "$$fmt_out" ]; then echo "Files need formatting:"; echo "$$fmt_out"; exit 1; fi
	@if command -v staticcheck >/dev/null 2>&1; then \
	  echo "Running staticcheck"; staticcheck ./...; \
	else \
	  echo "staticcheck not found (optional)"; \
	fi

.PHONY: test
test: ## Run tests (unit + integration) with race detector
	go test -race -count=1 ./...

.PHONY: cover
cover: ## Run tests with coverage profile (text summary)
	go test -race -covermode=atomic -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | tee coverage.func.txt
	@total=$$(go tool cover -func=coverage.out | grep '^total:' | awk '{print $$3}' | sed 's/%//'); \
	  awk 'BEGIN { if ('"$$total"' < $(MIN_COVER)) { printf "Coverage %.2f%% is below MIN_COVER=%.2f%%\n", '"$$total"', $(MIN_COVER); exit 1 } else { printf "Coverage %.2f%% (min %.2f%%) OK\n", '"$$total"', $(MIN_COVER) } }'

.PHONY: cover-html
cover-html: cover ## Generate HTML coverage report
	go tool cover -html=coverage.out -o coverage.html
	@echo "Open coverage.html in your browser."

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(BIN_DIR) $(DIST_DIR) coverage.out coverage.html coverage.func.txt

.PHONY: dist-clean
dist-clean: clean ## Alias for clean (legacy)

.PHONY: release
release: ## Build release archives for multiple OS/ARCH into dist/
	@echo "Building release artifacts for version $(VERSION)"
	@mkdir -p $(DIST_DIR)
	@rm -f $(DIST_DIR)/sha256sums.txt
	@set -e; \
	  for target in $(OS_ARCHES); do \
	    os=$${target%/*}; arch=$${target##*/}; \
	    ext=""; if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
	    out="$(DIST_DIR)/$(BINARY)-$(VERSION)-$$os-$$arch$$ext"; \
	    echo "  -> $$out"; \
	    GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o "$$out" $(PKG); \
	    sha256sum "$$out" >> $(DIST_DIR)/sha256sums.txt; \
	  done; \
	  sort -o $(DIST_DIR)/sha256sums.txt $(DIST_DIR)/sha256sums.txt
	@echo "Artifacts in $(DIST_DIR)/"
	@echo "SHA256 sums:"; cat $(DIST_DIR)/sha256sums.txt

.PHONY: release-check
release-check: ## Validate git state before tagging (no uncommitted changes, tag unused)
	@test -z "$$(git status --porcelain)" || (echo "Working tree not clean"; git status --short; exit 1)
	@if git rev-parse "$(TAG_NAME)" >/dev/null 2>&1; then echo "Tag $(TAG_NAME) already exists"; exit 1; fi
	@echo "Git state OK for version $(VERSION)"

.PHONY: tag
tag: release-check ## Create and push git tag $(TAG_NAME)
	git tag -a "$(TAG_NAME)" -m "Release $(VERSION)"
	git push origin "$(TAG_NAME)"

.PHONY: print-version
print-version: ## Print detected version
	@echo $(VERSION)

.PHONY: ci
ci: tidy lint test cover ## Run typical CI pipeline
	@echo "CI pipeline complete."

# Self-documentation: any target with '## desc' will show in 'make help'
