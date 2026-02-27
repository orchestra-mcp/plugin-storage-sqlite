.PHONY: proto build build-orchestrator build-storage-markdown build-tools-features build-transport-stdio test test-unit test-e2e clean

# Directories
BIN_DIR := bin
PROTO_DIR := proto
GEN_DIR := gen/go

# Binaries
ORCHESTRATOR := $(BIN_DIR)/orchestrator
STORAGE_MARKDOWN := $(BIN_DIR)/storage-markdown
TOOLS_FEATURES := $(BIN_DIR)/tools-features
TRANSPORT_STDIO := $(BIN_DIR)/transport-stdio

# === Proto ===

proto:
	cd $(PROTO_DIR) && buf lint && buf generate

# === Build ===

build: build-orchestrator build-storage-markdown build-tools-features build-transport-stdio

build-orchestrator:
	@mkdir -p $(BIN_DIR)
	go build -o $(ORCHESTRATOR) ./services/orchestrator/cmd/

build-storage-markdown:
	@mkdir -p $(BIN_DIR)
	go build -o $(STORAGE_MARKDOWN) ./plugins/storage-markdown/cmd/

build-tools-features:
	@mkdir -p $(BIN_DIR)
	go build -o $(TOOLS_FEATURES) ./plugins/tools-features/cmd/

build-transport-stdio:
	@mkdir -p $(BIN_DIR)
	go build -o $(TRANSPORT_STDIO) ./plugins/transport-stdio/cmd/

# === Test ===

test: test-unit

test-unit:
	go test ./libs/go/... -v
	go test ./services/orchestrator/... -v
	go test ./plugins/storage-markdown/... -v
	go test ./plugins/tools-features/... -v
	go test ./plugins/transport-stdio/... -v

test-e2e: build
	@bash scripts/test-e2e.sh

# === Clean ===

clean:
	rm -rf $(BIN_DIR)
	rm -rf ~/.orchestra/certs
