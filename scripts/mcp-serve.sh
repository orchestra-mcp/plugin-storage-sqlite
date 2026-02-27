#!/usr/bin/env bash
#
# Single-command MCP server for Claude Code / Cursor / etc.
# Starts orchestrator + plugins in background, then runs transport-stdio.
#
# Usage in .mcp.json:
#   { "mcpServers": { "orchestra": { "command": "./scripts/mcp-serve.sh" } } }
#
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BIN="$ROOT/bin"
WORKSPACE="${ORCHESTRA_WORKSPACE:-$ROOT}"
CERTS_DIR="${ORCHESTRA_CERTS_DIR:-$HOME/.orchestra/certs}"
LOG="${ORCHESTRA_LOG:-$ROOT/.orchestra-mcp.log}"
PIDFILE="$ROOT/.orchestra-mcp.pid"

# Build if binaries don't exist.
if [[ ! -f "$BIN/orchestrator" ]] || [[ ! -f "$BIN/transport-stdio" ]]; then
    make -C "$ROOT" build >/dev/null 2>&1
fi

# Kill ALL orchestra processes aggressively — no PID games.
pkill -9 -f "$BIN/orchestrator" 2>/dev/null || true
pkill -9 -f "$BIN/storage-markdown" 2>/dev/null || true
pkill -9 -f "$BIN/tools-features" 2>/dev/null || true
pkill -9 -f "$BIN/transport-stdio" 2>/dev/null || true
sleep 0.5

# Write a fresh config with port 0 (OS picks a free port).
TMPCONFIG=$(mktemp)

cat > "$TMPCONFIG" <<YAML
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

# Truncate the log so we only read fresh entries.
: > "$LOG"

# Cleanup: kill orchestrator + children, remove temp files.
ORCH_PID=""
cleanup() {
    if [[ -n "$ORCH_PID" ]]; then
        # Kill children (storage-markdown, tools-features) then orchestrator.
        pkill -P "$ORCH_PID" 2>/dev/null || true
        kill "$ORCH_PID" 2>/dev/null || true
        sleep 0.3
        pkill -9 -P "$ORCH_PID" 2>/dev/null || true
        kill -9 "$ORCH_PID" 2>/dev/null || true
    fi
    rm -f "$TMPCONFIG" "$PIDFILE"
}
trap cleanup EXIT INT TERM

# Start orchestrator in background.
"$BIN/orchestrator" --config "$TMPCONFIG" >>"$LOG" 2>&1 &
ORCH_PID=$!
echo "$ORCH_PID" > "$PIDFILE"

# Wait for both plugins to register.
READY=false
for i in $(seq 1 30); do
    sleep 0.5
    BOOTED=$(grep -c "registered and booted" "$LOG" 2>/dev/null || echo 0)
    if [[ "$BOOTED" -ge 2 ]]; then
        READY=true
        break
    fi
    if ! kill -0 "$ORCH_PID" 2>/dev/null; then
        echo "Orchestrator failed to start. Check $LOG" >&2
        exit 1
    fi
done

if [[ "$READY" != "true" ]]; then
    echo "Orchestrator did not become ready in 15 seconds. Check $LOG" >&2
    exit 1
fi

# Extract the actual listen address from the fresh log.
ORCH_ADDR=$(grep -o 'listening on [^ ]*' "$LOG" | head -1 | awk '{print $NF}')
if [[ -z "$ORCH_ADDR" ]]; then
    echo "Could not determine orchestrator address. Check $LOG" >&2
    exit 1
fi

# Run transport-stdio. NOT exec — trap must fire to clean up orchestrator.
"$BIN/transport-stdio" --orchestrator-addr="$ORCH_ADDR" --certs-dir="$CERTS_DIR"
