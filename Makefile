.PHONY: proto build build-orchestrator build-storage-markdown build-tools-features build-transport-stdio build-cli test test-unit test-e2e clean install release build-tools-marketplace

# Directories
BIN_DIR := bin
PROTO_DIR := libs/proto
GEN_DIR := libs/gen-go

# Binaries
ORCHESTRATOR := $(BIN_DIR)/orchestrator
STORAGE_MARKDOWN := $(BIN_DIR)/storage-markdown
TOOLS_FEATURES := $(BIN_DIR)/tools-features
TRANSPORT_STDIO := $(BIN_DIR)/transport-stdio
ORCHESTRA_CLI := $(BIN_DIR)/orchestra
TOOLS_MARKETPLACE := $(BIN_DIR)/tools-marketplace

# === Proto ===

proto:
	cd $(PROTO_DIR) && buf lint && buf generate

# === Build ===

build: build-orchestrator build-storage-markdown build-tools-features build-transport-stdio build-cli build-tools-marketplace

build-orchestrator:
	@mkdir -p $(BIN_DIR)
	go build -o $(ORCHESTRATOR) ./libs/orchestrator/cmd/

build-storage-markdown:
	@mkdir -p $(BIN_DIR)
	go build -o $(STORAGE_MARKDOWN) ./libs/plugin-storage-markdown/cmd/

build-tools-features:
	@mkdir -p $(BIN_DIR)
	go build -o $(TOOLS_FEATURES) ./libs/plugin-tools-features/cmd/

build-transport-stdio:
	@mkdir -p $(BIN_DIR)
	go build -o $(TRANSPORT_STDIO) ./libs/plugin-transport-stdio/cmd/

build-cli:
	@mkdir -p $(BIN_DIR)
	go build -ldflags "-X github.com/orchestra-mcp/cli/internal.Version=$(shell git describe --tags --always --dirty 2>/dev/null || echo dev) -X github.com/orchestra-mcp/cli/internal.Commit=$(shell git rev-parse --short HEAD 2>/dev/null || echo none) -X github.com/orchestra-mcp/cli/internal.Date=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -o $(ORCHESTRA_CLI) ./libs/cli/

# === Test ===

test: test-unit

test-unit:
	go test ./libs/sdk-go/... -v
	go test ./libs/orchestrator/... -v
	go test ./libs/plugin-storage-markdown/... -v
	go test ./libs/plugin-tools-features/... -v
	go test ./libs/plugin-transport-stdio/... -v
	go test ./libs/plugin-tools-marketplace/... -v

test-e2e: build
	@bash scripts/test-e2e.sh

# === Install ===

PREFIX ?= /usr/local
BINARIES := orchestra orchestrator storage-markdown tools-features transport-stdio tools-marketplace

install: build
	@mkdir -p $(PREFIX)/bin
	@for b in $(BINARIES); do \
		cp $(BIN_DIR)/$$b $(PREFIX)/bin/ && echo "  installed $(PREFIX)/bin/$$b"; \
	done
	@echo "\nOrchestra installed. Run: orchestra init"

uninstall:
	@for b in $(BINARIES); do \
		rm -f $(PREFIX)/bin/$$b && echo "  removed $(PREFIX)/bin/$$b"; \
	done

# === Release ===

DIST_DIR := dist
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X github.com/orchestra-mcp/cli/internal.Version=$(VERSION) \
           -X github.com/orchestra-mcp/cli/internal.Commit=$(shell git rev-parse --short HEAD 2>/dev/null || echo none) \
           -X github.com/orchestra-mcp/cli/internal.Date=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)

release:
	@mkdir -p $(DIST_DIR)
	@for goos in darwin linux; do \
		for goarch in amd64 arm64; do \
			echo "Building $$goos/$$goarch..."; \
			GOOS=$$goos GOARCH=$$goarch go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/orchestra ./libs/cli/; \
			GOOS=$$goos GOARCH=$$goarch go build -o $(DIST_DIR)/orchestrator ./libs/orchestrator/cmd/; \
			GOOS=$$goos GOARCH=$$goarch go build -o $(DIST_DIR)/storage-markdown ./libs/plugin-storage-markdown/cmd/; \
			GOOS=$$goos GOARCH=$$goarch go build -o $(DIST_DIR)/tools-features ./libs/plugin-tools-features/cmd/; \
			GOOS=$$goos GOARCH=$$goarch go build -o $(DIST_DIR)/transport-stdio ./libs/plugin-transport-stdio/cmd/; \
			GOOS=$$goos GOARCH=$$goarch go build -o $(DIST_DIR)/tools-marketplace ./libs/plugin-tools-marketplace/cmd/; \
			tar -czf $(DIST_DIR)/orchestra-$$goos-$$goarch.tar.gz -C $(DIST_DIR) orchestra orchestrator storage-markdown tools-features transport-stdio tools-marketplace; \
			rm -f $(DIST_DIR)/orchestra $(DIST_DIR)/orchestrator $(DIST_DIR)/storage-markdown $(DIST_DIR)/tools-features $(DIST_DIR)/transport-stdio $(DIST_DIR)/tools-marketplace; \
		done; \
	done
	@echo "\nRelease tarballs in $(DIST_DIR)/"
	@ls -lh $(DIST_DIR)/*.tar.gz

# === Clean ===

clean:
	rm -rf $(BIN_DIR) $(DIST_DIR)
	rm -rf ~/.orchestra/certs

build-tools-marketplace:
	@mkdir -p $(BIN_DIR)
	go build -o $(TOOLS_MARKETPLACE) ./libs/plugin-tools-marketplace/cmd/
