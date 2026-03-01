#!/usr/bin/env bash
# init-apps.sh — Push apps/ to their GitHub repos.

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

GREEN='\033[0;32m'; RED='\033[0;31m'; CYAN='\033[0;36m'; NC='\033[0m'
info()  { printf "${CYAN}[init]${NC}  %s\n" "$*"; }
ok()    { printf "${GREEN}[ok]${NC}    %s\n" "$*"; }
fail()  { printf "${RED}[fail]${NC}  %s\n" "$*" >&2; }

# App name | local path | GitHub URL
APPS=(
    "web|apps/web|https://github.com/orchestra-mcp/web.git"
    "next|apps/next|https://github.com/orchestra-mcp/next.git"
    "swift|apps/swift|https://github.com/orchestra-mcp/swift.git"
    "kotlin|apps/kotlin|https://github.com/orchestra-mcp/kotlin.git"
)

SUCCESS=0
ERRORS=0

for entry in "${APPS[@]}"; do
    IFS='|' read -r name path url <<< "$entry"
    src="${ROOT}/${path}"
    tmp="$(mktemp -d)"

    info "${name}: initializing..."

    git init --quiet "${tmp}/repo"
    git -C "${tmp}/repo" config user.name "Fady Mondy"
    git -C "${tmp}/repo" config user.email "info@3x1.io"

    # Copy source files, excluding build artifacts
    rsync -a \
        --exclude '.git' \
        --exclude 'node_modules' \
        --exclude '.next' \
        --exclude 'build' \
        --exclude 'DerivedData' \
        --exclude '.build' \
        --exclude '.gradle' \
        --exclude '*.xcuserdata' \
        --exclude 'target' \
        --exclude '.DS_Store' \
        "${src}/" "${tmp}/repo/"

    git -C "${tmp}/repo" add -A
    git -C "${tmp}/repo" commit -m "Initial commit" --quiet
    git -C "${tmp}/repo" branch -M master
    git -C "${tmp}/repo" remote add origin "${url}"

    if git -C "${tmp}/repo" push -u origin master --quiet 2>&1; then
        ok "${name}: pushed"
        SUCCESS=$((SUCCESS + 1))
    else
        fail "${name}: push failed"
        ERRORS=$((ERRORS + 1))
    fi

    rm -rf "$tmp"
done

echo ""
echo "────────────────────────────────"
ok "Initialized: ${SUCCESS}  Failed: ${ERRORS}"
[[ $ERRORS -gt 0 ]] && exit 1
exit 0
