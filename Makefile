.PHONY: proto build build-orchestrator build-storage-markdown build-tools-features build-transport-stdio build-cli build-web build-next dev-next storybook-next test test-unit test-e2e test-engine-rag clean install release build-tools-marketplace build-engine-rag build-bridge-claude build-tools-agentops build-tools-sessions build-tools-workspace build-transport-quic-bridge build-bridge-openai build-bridge-gemini build-bridge-ollama build-bridge-firecrawl build-tools-markdown build-tools-docs build-tools-notes build-devtools-git build-agent-orchestrator build-ai-screenshot build-ai-vision build-ai-browser-context build-ai-screen-reader build-services-voice build-services-notifications build-tools-extension-generator build-devtools-file-explorer build-devtools-terminal build-devtools-ssh build-devtools-services build-devtools-docker build-devtools-debugger build-devtools-test-runner build-devtools-log-viewer build-devtools-database build-devtools-devops build-integration-figma build-devtools-components xcodegen-swift build-swift build-swift-ios run-swift test-swift dev-swift clean-swift

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
DEVTOOLS_COMPONENTS := $(BIN_DIR)/devtools-components
INTEGRATION_FIGMA := $(BIN_DIR)/integration-figma
DEVTOOLS_DEVOPS := $(BIN_DIR)/devtools-devops
DEVTOOLS_DATABASE := $(BIN_DIR)/devtools-database
DEVTOOLS_LOG_VIEWER := $(BIN_DIR)/devtools-log-viewer
DEVTOOLS_TEST_RUNNER := $(BIN_DIR)/devtools-test-runner
DEVTOOLS_DEBUGGER := $(BIN_DIR)/devtools-debugger
DEVTOOLS_DOCKER := $(BIN_DIR)/devtools-docker
DEVTOOLS_SERVICES := $(BIN_DIR)/devtools-services
DEVTOOLS_SSH := $(BIN_DIR)/devtools-ssh
DEVTOOLS_TERMINAL := $(BIN_DIR)/devtools-terminal
DEVTOOLS_FILE_EXPLORER := $(BIN_DIR)/devtools-file-explorer
TOOLS_EXTENSION_GENERATOR := $(BIN_DIR)/tools-extension-generator
SERVICES_NOTIFICATIONS := $(BIN_DIR)/services-notifications
SERVICES_VOICE := $(BIN_DIR)/services-voice
AI_SCREEN_READER := $(BIN_DIR)/ai-screen-reader
AI_BROWSER_CONTEXT := $(BIN_DIR)/ai-browser-context
AI_VISION := $(BIN_DIR)/ai-vision
AI_SCREENSHOT := $(BIN_DIR)/ai-screenshot
AGENT_ORCHESTRATOR := $(BIN_DIR)/agent-orchestrator
DEVTOOLS_GIT := $(BIN_DIR)/devtools-git
TOOLS_NOTES := $(BIN_DIR)/tools-notes
TOOLS_DOCS := $(BIN_DIR)/tools-docs
TOOLS_MARKDOWN := $(BIN_DIR)/tools-markdown
BRIDGE_FIRECRAWL := $(BIN_DIR)/bridge-firecrawl
BRIDGE_OLLAMA := $(BIN_DIR)/bridge-ollama
BRIDGE_GEMINI := $(BIN_DIR)/bridge-gemini
BRIDGE_OPENAI := $(BIN_DIR)/bridge-openai
TRANSPORT_QUIC_BRIDGE := $(BIN_DIR)/transport-quic-bridge
TOOLS_SESSIONS := $(BIN_DIR)/tools-sessions
TOOLS_AGENTOPS := $(BIN_DIR)/tools-agentops
TOOLS_WORKSPACE := $(BIN_DIR)/tools-workspace
BRIDGE_CLAUDE := $(BIN_DIR)/bridge-claude
ENGINE_RAG := $(BIN_DIR)/engine-rag
TOOLS_MARKETPLACE := $(BIN_DIR)/tools-marketplace
WEB := $(BIN_DIR)/web

# === Proto ===

proto:
	cd $(PROTO_DIR) && buf lint && buf generate

# === Build ===

build: build-orchestrator build-storage-markdown build-tools-features build-transport-stdio build-cli build-web build-tools-marketplace build-engine-rag build-bridge-claude build-tools-agentops build-tools-sessions build-tools-workspace build-transport-quic-bridge build-bridge-openai build-bridge-gemini build-bridge-ollama build-bridge-firecrawl build-tools-markdown build-tools-docs build-tools-notes build-devtools-git build-agent-orchestrator build-ai-screenshot build-ai-vision build-ai-browser-context build-ai-screen-reader build-services-voice build-services-notifications build-tools-extension-generator build-devtools-file-explorer build-devtools-terminal build-devtools-ssh build-devtools-services build-devtools-docker build-devtools-debugger build-devtools-test-runner build-devtools-log-viewer build-devtools-database build-devtools-devops build-integration-figma build-devtools-components

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

build-web:
	@mkdir -p $(BIN_DIR)
	go build -o $(WEB) ./apps/web/cmd/

NEXT_DIR := apps/next

build-next:
	cd $(NEXT_DIR) && npm install && npm run build

dev-next:
	cd $(NEXT_DIR) && npm run dev

storybook-next:
	cd $(NEXT_DIR) && npm run storybook

# === Test ===

test: test-unit

test-unit:
	go test ./libs/sdk-go/... -v
	go test ./libs/orchestrator/... -v
	go test ./libs/plugin-storage-markdown/... -v
	go test ./libs/plugin-tools-features/... -v
	go test ./libs/plugin-transport-stdio/... -v
	go test ./libs/plugin-devtools-components/... -v
	go test ./libs/plugin-integration-figma/... -v
	go test ./libs/plugin-devtools-devops/... -v
	go test ./libs/plugin-devtools-database/... -v
	go test ./libs/plugin-devtools-log-viewer/... -v
	go test ./libs/plugin-devtools-test-runner/... -v
	go test ./libs/plugin-devtools-debugger/... -v
	go test ./libs/plugin-devtools-docker/... -v
	go test ./libs/plugin-devtools-services/... -v
	go test ./libs/plugin-devtools-ssh/... -v
	go test ./libs/plugin-devtools-terminal/... -v
	go test ./libs/plugin-devtools-file-explorer/... -v
	go test ./libs/plugin-tools-extension-generator/... -v
	go test ./libs/plugin-services-notifications/... -v
	go test ./libs/plugin-services-voice/... -v
	go test ./libs/plugin-ai-screen-reader/... -v
	go test ./libs/plugin-ai-browser-context/... -v
	go test ./libs/plugin-ai-vision/... -v
	go test ./libs/plugin-ai-screenshot/... -v
	go test ./libs/plugin-agent-orchestrator/... -v
	go test ./libs/plugin-devtools-git/... -v
	go test ./libs/plugin-tools-notes/... -v
	go test ./libs/plugin-tools-docs/... -v
	go test ./libs/plugin-tools-markdown/... -v
	go test ./libs/plugin-bridge-firecrawl/... -v
	go test ./libs/plugin-bridge-ollama/... -v
	go test ./libs/plugin-bridge-gemini/... -v
	go test ./libs/plugin-bridge-openai/... -v
	go test ./libs/plugin-transport-quic-bridge/... -v
	go test ./libs/plugin-tools-sessions/... -v
	go test ./libs/plugin-tools-agentops/... -v
	go test ./libs/plugin-tools-workspace/... -v
	go test ./libs/plugin-bridge-claude/... -v
	go test ./libs/plugin-tools-marketplace/... -v

test-e2e: build
	@bash scripts/test-e2e.sh

# === Install ===

PREFIX ?= /usr/local
BINARIES := orchestra orchestrator storage-markdown tools-features transport-stdio web tools-marketplace engine-rag bridge-claude tools-agentops tools-sessions tools-workspace transport-quic-bridge bridge-openai bridge-gemini bridge-ollama bridge-firecrawl tools-markdown tools-docs tools-notes devtools-git agent-orchestrator ai-screenshot ai-vision ai-browser-context ai-screen-reader services-voice services-notifications tools-extension-generator devtools-file-explorer devtools-terminal devtools-ssh devtools-services devtools-docker devtools-debugger devtools-test-runner devtools-log-viewer devtools-database devtools-devops integration-figma devtools-components

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

build-engine-rag:
	@mkdir -p $(BIN_DIR)
	cd libs/plugin-engine-rag && cargo build --release
	cp libs/plugin-engine-rag/target/release/orchestra-rag $(ENGINE_RAG)

test-engine-rag:
	cd libs/plugin-engine-rag && cargo test

build-bridge-claude:
	@mkdir -p $(BIN_DIR)
	go build -o $(BRIDGE_CLAUDE) ./libs/plugin-bridge-claude/cmd/

build-tools-agentops:
	@mkdir -p $(BIN_DIR)
	go build -o $(TOOLS_AGENTOPS) ./libs/plugin-tools-agentops/cmd/

build-tools-sessions:
	@mkdir -p $(BIN_DIR)
	go build -o $(TOOLS_SESSIONS) ./libs/plugin-tools-sessions/cmd/

build-tools-workspace:
	@mkdir -p $(BIN_DIR)
	go build -o $(TOOLS_WORKSPACE) ./libs/plugin-tools-workspace/cmd/

build-transport-quic-bridge:
	@mkdir -p $(BIN_DIR)
	go build -o $(TRANSPORT_QUIC_BRIDGE) ./libs/plugin-transport-quic-bridge/cmd/

build-bridge-openai:
	@mkdir -p $(BIN_DIR)
	go build -o $(BRIDGE_OPENAI) ./libs/plugin-bridge-openai/cmd/

build-bridge-gemini:
	@mkdir -p $(BIN_DIR)
	go build -o $(BRIDGE_GEMINI) ./libs/plugin-bridge-gemini/cmd/

build-bridge-ollama:
	@mkdir -p $(BIN_DIR)
	go build -o $(BRIDGE_OLLAMA) ./libs/plugin-bridge-ollama/cmd/

build-bridge-firecrawl:
	@mkdir -p $(BIN_DIR)
	go build -o $(BRIDGE_FIRECRAWL) ./libs/plugin-bridge-firecrawl/cmd/

build-tools-markdown:
	@mkdir -p $(BIN_DIR)
	go build -o $(TOOLS_MARKDOWN) ./libs/plugin-tools-markdown/cmd/

build-tools-docs:
	@mkdir -p $(BIN_DIR)
	go build -o $(TOOLS_DOCS) ./libs/plugin-tools-docs/cmd/

build-tools-notes:
	@mkdir -p $(BIN_DIR)
	go build -o $(TOOLS_NOTES) ./libs/plugin-tools-notes/cmd/

build-devtools-git:
	@mkdir -p $(BIN_DIR)
	go build -o $(DEVTOOLS_GIT) ./libs/plugin-devtools-git/cmd/

build-agent-orchestrator:
	@mkdir -p $(BIN_DIR)
	go build -o $(AGENT_ORCHESTRATOR) ./libs/plugin-agent-orchestrator/cmd/

build-ai-screenshot:
	@mkdir -p $(BIN_DIR)
	go build -o $(AI_SCREENSHOT) ./libs/plugin-ai-screenshot/cmd/

build-ai-vision:
	@mkdir -p $(BIN_DIR)
	go build -o $(AI_VISION) ./libs/plugin-ai-vision/cmd/

build-ai-browser-context:
	@mkdir -p $(BIN_DIR)
	go build -o $(AI_BROWSER_CONTEXT) ./libs/plugin-ai-browser-context/cmd/

build-ai-screen-reader:
	@mkdir -p $(BIN_DIR)
	go build -o $(AI_SCREEN_READER) ./libs/plugin-ai-screen-reader/cmd/

build-services-voice:
	@mkdir -p $(BIN_DIR)
	go build -o $(SERVICES_VOICE) ./libs/plugin-services-voice/cmd/

build-services-notifications:
	@mkdir -p $(BIN_DIR)
	go build -o $(SERVICES_NOTIFICATIONS) ./libs/plugin-services-notifications/cmd/

build-tools-extension-generator:
	@mkdir -p $(BIN_DIR)
	go build -o $(TOOLS_EXTENSION_GENERATOR) ./libs/plugin-tools-extension-generator/cmd/

build-devtools-file-explorer:
	@mkdir -p $(BIN_DIR)
	go build -o $(DEVTOOLS_FILE_EXPLORER) ./libs/plugin-devtools-file-explorer/cmd/

build-devtools-terminal:
	@mkdir -p $(BIN_DIR)
	go build -o $(DEVTOOLS_TERMINAL) ./libs/plugin-devtools-terminal/cmd/

build-devtools-ssh:
	@mkdir -p $(BIN_DIR)
	go build -o $(DEVTOOLS_SSH) ./libs/plugin-devtools-ssh/cmd/

build-devtools-services:
	@mkdir -p $(BIN_DIR)
	go build -o $(DEVTOOLS_SERVICES) ./libs/plugin-devtools-services/cmd/

build-devtools-docker:
	@mkdir -p $(BIN_DIR)
	go build -o $(DEVTOOLS_DOCKER) ./libs/plugin-devtools-docker/cmd/

build-devtools-debugger:
	@mkdir -p $(BIN_DIR)
	go build -o $(DEVTOOLS_DEBUGGER) ./libs/plugin-devtools-debugger/cmd/

build-devtools-test-runner:
	@mkdir -p $(BIN_DIR)
	go build -o $(DEVTOOLS_TEST_RUNNER) ./libs/plugin-devtools-test-runner/cmd/

build-devtools-log-viewer:
	@mkdir -p $(BIN_DIR)
	go build -o $(DEVTOOLS_LOG_VIEWER) ./libs/plugin-devtools-log-viewer/cmd/

build-devtools-database:
	@mkdir -p $(BIN_DIR)
	go build -o $(DEVTOOLS_DATABASE) ./libs/plugin-devtools-database/cmd/

build-devtools-devops:
	@mkdir -p $(BIN_DIR)
	go build -o $(DEVTOOLS_DEVOPS) ./libs/plugin-devtools-devops/cmd/

build-integration-figma:
	@mkdir -p $(BIN_DIR)
	go build -o $(INTEGRATION_FIGMA) ./libs/plugin-integration-figma/cmd/

build-devtools-components:
	@mkdir -p $(BIN_DIR)
	go build -o $(DEVTOOLS_COMPONENTS) ./libs/plugin-devtools-components/cmd/

# === Swift Universal App ===

SWIFT_DIR        := apps/swift
SWIFT_SCHEME_MAC := OrchestraMac
SWIFT_SCHEME_IOS := OrchestraIOS
SWIFT_CONFIG     := Debug
SWIFT_XCPROJ     := $(SWIFT_DIR)/Orchestra.xcodeproj

# Derived data app path helper (resolves after build)
SWIFT_APP = $(shell find ~/Library/Developer/Xcode/DerivedData/Orchestra-*/Build/Products/Debug \
	-name "Orchestra.app" -not -path "*/Index.noindex/*" 2>/dev/null | head -1)

.PHONY: xcodegen-swift build-swift run-swift test-swift dev-swift clean-swift

xcodegen-swift: ## Regenerate Orchestra.xcodeproj from project.yml
	xcodegen generate --spec $(SWIFT_DIR)/project.yml

build-swift: xcodegen-swift ## Build macOS Swift app (Debug)
	xcodebuild \
	  -project $(SWIFT_XCPROJ) \
	  -scheme $(SWIFT_SCHEME_MAC) \
	  -configuration $(SWIFT_CONFIG) \
	  -destination "platform=macOS" \
	  build | xcpretty 2>/dev/null || xcodebuild \
	  -project $(SWIFT_XCPROJ) \
	  -scheme $(SWIFT_SCHEME_MAC) \
	  -configuration $(SWIFT_CONFIG) \
	  -destination "platform=macOS" \
	  build

build-swift-ios: xcodegen-swift ## Build iOS Swift app (simulator)
	xcodebuild \
	  -project $(SWIFT_XCPROJ) \
	  -scheme $(SWIFT_SCHEME_IOS) \
	  -configuration $(SWIFT_CONFIG) \
	  -destination "platform=iOS Simulator,name=iPhone 16" \
	  build | xcpretty 2>/dev/null || true

run-swift: build-swift ## Build and launch Orchestra.app on macOS
	@APP="$(SWIFT_APP)"; \
	if [ -z "$$APP" ]; then echo "Orchestra.app not found — run make build-swift first"; exit 1; fi; \
	pkill -x Orchestra 2>/dev/null; sleep 0.3; \
	open "$$APP"

test-swift: ## Run Swift package unit tests (SPM)
	cd $(SWIFT_DIR) && swift test

dev-swift: ## Watch Swift sources and auto-rebuild + relaunch on changes (requires fswatch)
	@which fswatch > /dev/null 2>&1 || { echo "Installing fswatch..."; brew install fswatch; }
	@echo "Watching $(SWIFT_DIR) for changes (Ctrl+C to stop)..."
	@echo "Initial build..."
	@$(MAKE) build-swift && $(MAKE) run-swift || true
	@fswatch -o $(SWIFT_DIR) \
	  --exclude ".*\.xcodeproj" \
	  --exclude ".*\.build" \
	  --exclude ".*DerivedData" \
	  --exclude ".*\.resolved" \
	  | xargs -n1 -I{} sh -c ' \
	    echo "\n[dev-swift] Change detected — rebuilding..."; \
	    xcodebuild \
	      -project $(SWIFT_XCPROJ) \
	      -scheme $(SWIFT_SCHEME_MAC) \
	      -configuration $(SWIFT_CONFIG) \
	      -destination "platform=macOS" \
	      build -quiet 2>&1 | tail -5 && \
	    APP=$$(find ~/Library/Developer/Xcode/DerivedData/Orchestra-*/Build/Products/Debug \
	      -name "Orchestra.app" -not -path "*/Index.noindex/*" 2>/dev/null | head -1); \
	    pkill -x Orchestra 2>/dev/null; sleep 0.3; \
	    open "$$APP" && echo "[dev-swift] Relaunched Orchestra.app" \
	    || echo "[dev-swift] Build failed — fix errors and save again" \
	  '

clean-swift: ## Remove Swift DerivedData for Orchestra
	rm -rf ~/Library/Developer/Xcode/DerivedData/Orchestra-*
