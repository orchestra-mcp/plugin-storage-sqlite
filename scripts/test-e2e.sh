#!/usr/bin/env bash
#
# End-to-end test for Phase 1 of the Orchestra Plugin Host.
#
# Flow: transport.stdio → orchestrator → tools.features → orchestrator → storage.markdown → disk
#
# Prerequisites: make build (all binaries in bin/)
#
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BIN="$ROOT/bin"

# Use a temporary workspace so we don't pollute the repo.
WORKSPACE="$(mktemp -d)"
CERTS_DIR="$WORKSPACE/.certs"

# PID tracking for cleanup.
ORCH_PID=""

cleanup() {
    echo "--- Cleaning up ---"
    if [[ -n "$ORCH_PID" ]]; then
        kill "$ORCH_PID" 2>/dev/null || true
        wait "$ORCH_PID" 2>/dev/null || true
    fi
    rm -rf "$WORKSPACE"
    echo "--- Done ---"
}
trap cleanup EXIT

echo "=== Orchestra E2E Test ==="
echo "Workspace: $WORKSPACE"
echo "Certs dir: $CERTS_DIR"

# Write a plugins.yaml for this test with absolute paths.
cat > "$WORKSPACE/plugins.yaml" <<YAML
listen_addr: "localhost:0"
certs_dir: "$CERTS_DIR"
plugins:
  - id: storage.markdown
    binary: $BIN/storage-markdown
    enabled: true
    provides_storage:
      - markdown
    args:
      - "--workspace=$WORKSPACE"
  - id: tools.features
    binary: $BIN/tools-features
    enabled: true
YAML

echo ""
echo "--- Step 1: Start orchestrator ---"

# Start orchestrator in background, capture stderr to read the listen address.
ORCH_LOG="$WORKSPACE/orchestrator.log"
"$BIN/orchestrator" --config "$WORKSPACE/plugins.yaml" > /dev/null 2>"$ORCH_LOG" &
ORCH_PID=$!

# Wait for orchestrator to be ready (it logs "started QUIC server" or plugin registrations).
echo "Waiting for orchestrator (pid=$ORCH_PID)..."
READY=false
for i in $(seq 1 30); do
    sleep 1
    if grep -q "registered and booted" "$ORCH_LOG" 2>/dev/null; then
        READY=true
        break
    fi
    # Check if the process is still alive.
    if ! kill -0 "$ORCH_PID" 2>/dev/null; then
        echo "FAIL: Orchestrator exited prematurely"
        cat "$ORCH_LOG"
        exit 1
    fi
done

if [[ "$READY" != "true" ]]; then
    echo "FAIL: Orchestrator did not become ready in 30 seconds"
    cat "$ORCH_LOG"
    exit 1
fi

echo "Orchestrator is ready."

# Extract the listen address from the orchestrator log.
# The orchestrator logs "QUIC server listening on <addr>"
ORCH_ADDR=$(grep -o 'listening on [^ ]*' "$ORCH_LOG" | head -1 | awk '{print $NF}')
if [[ -z "$ORCH_ADDR" ]]; then
    # Fallback: try to parse from config (localhost:0 means random port).
    # If we can't find it, use the default.
    echo "WARNING: Could not extract orchestrator address from log, trying to find it..."
    ORCH_ADDR=$(grep -o 'localhost:[0-9]*' "$ORCH_LOG" | tail -1)
fi

if [[ -z "$ORCH_ADDR" ]]; then
    echo "FAIL: Could not determine orchestrator listen address"
    cat "$ORCH_LOG"
    exit 1
fi

echo "Orchestrator listening at: $ORCH_ADDR"

echo ""
echo "--- Step 2: Send MCP initialize ---"

INIT_REQ='{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"e2e-test","version":"0.1.0"}}}'
INIT_RESP=$(echo "$INIT_REQ" | "$BIN/transport-stdio" --orchestrator-addr="$ORCH_ADDR" --certs-dir="$CERTS_DIR")

echo "Response: $INIT_RESP"

if echo "$INIT_RESP" | grep -q '"protocolVersion"'; then
    echo "PASS: initialize returned protocolVersion"
else
    echo "FAIL: initialize did not return protocolVersion"
    exit 1
fi

echo ""
echo "--- Step 3: Send tools/list ---"

LIST_REQ='{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'
LIST_RESP=$(echo "$LIST_REQ" | "$BIN/transport-stdio" --orchestrator-addr="$ORCH_ADDR" --certs-dir="$CERTS_DIR")

echo "Response (truncated): $(echo "$LIST_RESP" | head -c 500)"

if echo "$LIST_RESP" | grep -q '"create_project"'; then
    echo "PASS: tools/list includes create_project"
else
    echo "FAIL: tools/list does not include create_project"
    exit 1
fi

if echo "$LIST_RESP" | grep -q '"create_feature"'; then
    echo "PASS: tools/list includes create_feature"
else
    echo "FAIL: tools/list does not include create_feature"
    exit 1
fi

echo ""
echo "--- Step 4: Create a project ---"

CREATE_PROJECT_REQ='{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"create_project","arguments":{"name":"test-project","description":"E2E test project"}}}'
CREATE_PROJECT_RESP=$(echo "$CREATE_PROJECT_REQ" | "$BIN/transport-stdio" --orchestrator-addr="$ORCH_ADDR" --certs-dir="$CERTS_DIR")

echo "Response: $CREATE_PROJECT_RESP"

if echo "$CREATE_PROJECT_RESP" | grep -q 'test-project'; then
    echo "PASS: create_project returned project data"
else
    echo "FAIL: create_project did not return expected data"
    exit 1
fi

echo ""
echo "--- Step 5: Create a feature ---"

CREATE_FEATURE_REQ='{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"create_feature","arguments":{"project_id":"test-project","title":"Test Feature","description":"Created by E2E test","priority":"P1"}}}'
CREATE_FEATURE_RESP=$(echo "$CREATE_FEATURE_REQ" | "$BIN/transport-stdio" --orchestrator-addr="$ORCH_ADDR" --certs-dir="$CERTS_DIR")

echo "Response: $CREATE_FEATURE_RESP"

if echo "$CREATE_FEATURE_RESP" | grep -q 'FEAT-'; then
    echo "PASS: create_feature returned feature with FEAT- ID"
else
    echo "FAIL: create_feature did not return expected data"
    exit 1
fi

echo ""
echo "--- Step 6: Verify feature file on disk ---"

# List the features directory.
FEATURES_DIR="$WORKSPACE/.projects/test-project/features"
if [[ -d "$FEATURES_DIR" ]]; then
    echo "Features directory exists: $FEATURES_DIR"
    ls -la "$FEATURES_DIR/"

    # Check at least one FEAT-*.md file exists.
    FEAT_FILES=$(find "$FEATURES_DIR" -name "FEAT-*.md" 2>/dev/null | wc -l | tr -d ' ')
    if [[ "$FEAT_FILES" -gt 0 ]]; then
        echo "PASS: Found $FEAT_FILES feature file(s) on disk"

        # Show the content of the first feature file.
        FIRST_FEAT=$(find "$FEATURES_DIR" -name "FEAT-*.md" | head -1)
        echo ""
        echo "Feature file content:"
        cat "$FIRST_FEAT"
        echo ""

        # Verify it contains YAML frontmatter (--- delimiters with metadata).
        if grep -q "^---" "$FIRST_FEAT"; then
            echo "PASS: Feature file contains YAML frontmatter"
        else
            echo "FAIL: Feature file missing YAML frontmatter"
            exit 1
        fi
    else
        echo "FAIL: No FEAT-*.md files found in $FEATURES_DIR"
        exit 1
    fi
else
    echo "FAIL: Features directory does not exist at $FEATURES_DIR"
    ls -laR "$WORKSPACE/.projects/" 2>/dev/null || echo "No .projects directory"
    exit 1
fi

echo ""
echo "--- Step 7: Get feature back ---"

# Extract feature ID from the create response.
FEATURE_ID=$(echo "$CREATE_FEATURE_RESP" | grep -o 'FEAT-[A-Z0-9]*' | head -1)
if [[ -z "$FEATURE_ID" ]]; then
    echo "WARNING: Could not extract feature ID, skipping get_feature test"
else
    GET_FEATURE_REQ="{\"jsonrpc\":\"2.0\",\"id\":5,\"method\":\"tools/call\",\"params\":{\"name\":\"get_feature\",\"arguments\":{\"project_id\":\"test-project\",\"feature_id\":\"$FEATURE_ID\"}}}"
    GET_FEATURE_RESP=$(echo "$GET_FEATURE_REQ" | "$BIN/transport-stdio" --orchestrator-addr="$ORCH_ADDR" --certs-dir="$CERTS_DIR")

    echo "Response: $GET_FEATURE_RESP"

    if echo "$GET_FEATURE_RESP" | grep -q 'Test Feature'; then
        echo "PASS: get_feature returned correct title"
    else
        echo "FAIL: get_feature did not return expected data"
        exit 1
    fi
fi

echo ""
echo "==============================="
echo "  ALL E2E TESTS PASSED"
echo "==============================="
