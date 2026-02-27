# Step 8: CLI Binary + Distribution

## Goal

Replace the shell-script MCP launcher with a proper Go CLI binary (`orchestra`) that provides `orchestra init` (auto-configure any AI IDE) and `orchestra serve` (start MCP server), plus a full distribution pipeline.

## What Was Built

### `cmd/orchestra/` — New Go Module

Standalone module with minimal dependencies (`gopkg.in/yaml.v3` only). Does NOT import the plugin SDK, protobuf, or QUIC — uses `os/exec` to run sibling binaries.

**Files:**
- `main.go` — Subcommand dispatch (init, serve, version, help). Default: serve.
- `internal/serve.go` — Full Go replacement of `scripts/mcp-serve.sh`
- `internal/initcmd.go` — IDE config generation command
- `internal/ide.go` — 9 IDE config format definitions + generators
- `internal/detect.go` — Project name/type auto-detection
- `internal/version.go` — Version info via ldflags

### `orchestra serve` — MCP Server Launcher

Replaces `scripts/mcp-serve.sh` with equivalent Go code:

1. Resolves sibling binaries via `os.Executable()` + `filepath.EvalSymlinks()`
2. Kills stale processes via `pkill -9 -f`
3. Writes temp `plugins.yaml` (listen_addr: `localhost:0`)
4. Starts orchestrator subprocess, redirects stderr to log file
5. Polls log for `"registered and booted"` (2 plugins, 15s timeout)
6. Extracts actual listen address from log
7. Runs transport-stdio with `cmd.Stdin = os.Stdin`, `cmd.Stdout = os.Stdout`
8. Signal handler + cleanup (kill children, remove temp files)

**Critical constraint:** Never writes to stdout (it's MCP's JSON-RPC channel).

Default subcommand — running `orchestra` with no args starts the MCP server. This allows MCP clients to use `"command": "orchestra"` directly.

### `orchestra init` — IDE Config Generator

One command to configure MCP for any AI IDE:

```bash
orchestra init              # Auto-detect IDE from existing config dirs
orchestra init --ide=cursor # Configure specific IDE
orchestra init --all        # Generate all 9 IDE configs
```

**9 supported IDEs:**

| IDE | Config Path | Format |
|-----|-------------|--------|
| Claude Code | `.mcp.json` | JSON (`mcpServers`) |
| Cursor | `.cursor/mcp.json` | JSON (`mcpServers`) |
| VS Code / Copilot | `.vscode/mcp.json` | JSON (`mcpServers`) |
| Cline | `.vscode/mcp.json` | JSON (`mcpServers`) |
| Windsurf | `~/.codeium/windsurf/mcp_config.json` | JSON (`mcpServers`) — global |
| Codex (OpenAI) | `.codex/config.toml` | TOML (`mcp_servers`) |
| Gemini Code Assist | `.gemini/settings.json` | JSON (`mcpServers`) |
| Zed | `.zed/settings.json` | JSON (`context_servers`) |
| Continue.dev | `.continue/mcpServers/orchestra.yaml` | YAML |

**Key behaviors:**
- Reads existing config files and merges — never overwrites other MCP servers
- Creates parent directories as needed
- Creates `.projects/` directory for feature storage
- Detects project name from `package.json` / `go.mod` / `Cargo.toml` / `pyproject.toml`
- Auto-detects IDE from existing config dirs, defaults to Claude Code
- Codex TOML and Continue YAML generated via string templates (no extra deps)
- All output to stderr, never stdout

### Distribution Pipeline

**Makefile targets:**
- `make install` — Copy 5 binaries to `/usr/local/bin`
- `make uninstall` — Remove from `/usr/local/bin`
- `make release` — Cross-compile darwin/linux × amd64/arm64 → 4 tarballs in `dist/`

**`scripts/install.sh`** — curl | sh installer:
- Detects OS (darwin/linux) and arch (amd64/arm64)
- Downloads correct tarball from GitHub releases
- Installs all 5 binaries to `/usr/local/bin`

**GitHub Actions:**
- `.github/workflows/ci.yml` — Build + test + lint on push/PR
- `.github/workflows/release.yml` — On tag `v*`: cross-compile → package → GitHub Release

**Packaging templates:**
- `packaging/homebrew/orchestra.rb` — Homebrew formula (SHA256 placeholders)
- `packaging/npm/` — npm wrapper with postinstall binary download

## Design Decisions

### 1. Standalone Module
`cmd/orchestra` has no dependency on the plugin SDK or protobuf. It only uses `os/exec` + `yaml.v3`. This keeps the binary small and avoids pulling in QUIC/crypto deps.

### 2. Default to Serve
Running `orchestra` with no subcommand starts `serve`. MCP clients declare `"command": "orchestra"` without needing to know about subcommands.

### 3. Sibling Binary Discovery
`orchestra serve` finds `orchestrator`, `storage-markdown`, etc. in the same directory as itself via `os.Executable()`. All 5 binaries must be co-located. This works naturally with tarballs, Homebrew, and npm installs.

### 4. Config Merge
`orchestra init` reads existing JSON config files, adds/updates only the `orchestra` entry, and writes back. Other MCP servers are preserved. Running init twice is safe (idempotent).

### 5. No TOML/YAML Libraries for Simple Configs
Codex (TOML) and Continue.dev (YAML) configs are simple enough to generate with `fmt.Sprintf` templates. Only `yaml.v3` is imported for the orchestrator config.

## Verification

```bash
# Build
make build                    # 5 binaries in bin/

# CLI commands
bin/orchestra version         # Print version
bin/orchestra help            # Show usage

# Init
bin/orchestra init --all      # Generate 9 IDE configs
bin/orchestra init --ide=claude,cursor  # Specific IDEs

# Serve (test with MCP JSON-RPC)
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | bin/orchestra serve

# Install globally
make install                  # 5 binaries → /usr/local/bin

# Cross-compile release
make release                  # 4 tarballs in dist/

# Tests
make test                     # 66 unit tests pass
```
