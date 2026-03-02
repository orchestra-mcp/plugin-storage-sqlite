#!/usr/bin/env bash
# add-plugin-release-workflow.sh — Adds release.yml workflow to all optional plugin directories.
# This enables per-plugin GitHub releases with pre-built binaries for 4 platforms.
#
# Usage:
#   ./scripts/add-plugin-release-workflow.sh              # Add to all optional plugins
#   ./scripts/add-plugin-release-workflow.sh --dry-run     # Preview without writing
#   ./scripts/add-plugin-release-workflow.sh plugin-tools-notes  # Add to specific plugin

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DRY_RUN=false
FILTER=()

# Core plugins (bundled in-process, don't need release workflow)
CORE_PLUGINS="plugin-storage-markdown plugin-tools-features plugin-tools-marketplace plugin-transport-stdio"
# Non-plugin packages (libs, protos, CLI)
NON_PLUGINS="proto gen-go sdk-go orchestrator cli plugin-engine-rag"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --dry-run) DRY_RUN=true; shift ;;
    --help|-h)
      echo "Usage: $0 [--dry-run] [plugin-name ...]"
      exit 0
      ;;
    *) FILTER+=("$1"); shift ;;
  esac
done

CYAN='\033[0;36m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'
info()  { printf "${CYAN}[workflow]${NC}  %s\n" "$*"; }
ok()    { printf "${GREEN}[ok]${NC}       %s\n" "$*"; }
skip()  { printf "${YELLOW}[skip]${NC}     %s\n" "$*"; }

generate_release_yml() {
  local plugin_dir="$1"  # e.g. plugin-tools-notes
  local binary_name="$plugin_dir"  # Use full dir name as binary name (matches install.go)

  cat << 'WORKFLOW_EOF'
name: Release

on:
  push:
    tags: ['v*']

permissions:
  contents: write

env:
  GONOSUMCHECK: github.com/orchestra-mcp/*
  GONOSUMDB: github.com/orchestra-mcp/*
  GOFLAGS: -mod=mod

jobs:
  build:
    strategy:
      matrix:
        include:
          - goos: darwin
            goarch: amd64
            runner: macos-latest
          - goos: darwin
            goarch: arm64
            runner: macos-latest
          - goos: linux
            goarch: amd64
            runner: ubuntu-latest
          - goos: linux
            goarch: arm64
            runner: ubuntu-latest
    runs-on: ${{ matrix.runner }}
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Download dependencies
        run: |
          for i in 1 2 3; do
            go mod download && break
            echo "Retry $i: waiting for Go proxy to index..."
            sleep 15
          done

      - name: Build
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
          CGO_ENABLED: '0'
        run: |
WORKFLOW_EOF

  # This part needs the actual binary name interpolated
  echo "          BINARY_NAME=\"${binary_name}\""
  cat << 'WORKFLOW_EOF'
          go build -trimpath -ldflags "-s -w" -o "${BINARY_NAME}" ./cmd/

      - name: Package tarball
        run: |
WORKFLOW_EOF

  echo "          BINARY_NAME=\"${binary_name}\""
  cat << 'WORKFLOW_EOF'
          TARBALL="${BINARY_NAME}-${{ matrix.goos }}-${{ matrix.goarch }}.tar.gz"
          tar -czf "${TARBALL}" "${BINARY_NAME}"
          echo "TARBALL=${TARBALL}" >> $GITHUB_ENV

      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: ${{ env.TARBALL }}
          path: ${{ env.TARBALL }}

  release:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/download-artifact@v4
        with:
          path: artifacts
          merge-multiple: true

      - name: Create GitHub Release
        uses: softprops/action-gh-release@v2
        with:
          generate_release_notes: true
          files: artifacts/*.tar.gz
WORKFLOW_EOF
}

count=0
skipped=0

for dir in "${ROOT}"/libs/plugin-*; do
  [[ ! -d "$dir" ]] && continue
  plugin_name="$(basename "$dir")"

  # Skip core plugins
  is_core=false
  for core in $CORE_PLUGINS; do
    [[ "$plugin_name" == "$core" ]] && is_core=true
  done
  $is_core && { skip "${plugin_name} (core — bundled in-process)"; skipped=$((skipped + 1)); continue; }

  # Skip non-plugin packages
  is_non=false
  for non in $NON_PLUGINS; do
    [[ "$plugin_name" == "$non" ]] && is_non=true
  done
  $is_non && { skip "${plugin_name} (not a standalone plugin)"; skipped=$((skipped + 1)); continue; }

  # Apply filter if provided
  if [[ ${#FILTER[@]} -gt 0 ]]; then
    matched=false
    for f in "${FILTER[@]}"; do
      [[ "$plugin_name" == "$f" ]] && matched=true
    done
    $matched || continue
  fi

  workflow_dir="${dir}/.github/workflows"
  workflow_file="${workflow_dir}/release.yml"

  if $DRY_RUN; then
    info "${plugin_name}: would create ${workflow_file#"${ROOT}/"}"
    count=$((count + 1))
    continue
  fi

  mkdir -p "$workflow_dir"
  generate_release_yml "$plugin_name" > "$workflow_file"
  ok "${plugin_name}: created release.yml"
  count=$((count + 1))
done

echo ""
echo "────────────────────────────────"
echo -e "${GREEN}[done]${NC}  Added: ${count}  Skipped: ${skipped}"
if ! $DRY_RUN && [[ $count -gt 0 ]]; then
  echo ""
  info "Next: run ./scripts/sync-repos.sh to push workflows to GitHub"
  info "Then: ./scripts/release.sh v0.2.0 to tag all repos (triggers the workflows)"
fi
