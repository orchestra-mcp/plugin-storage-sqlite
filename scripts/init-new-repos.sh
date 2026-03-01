#!/usr/bin/env bash
# init-new-repos.sh — Push initial code to newly created empty GitHub repos.
# Only handles repos that don't have a master branch yet.

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOCK_FILE="${ROOT}/orchestra.lock"

# Repos that already existed before this run (have master branch + history)
EXISTING="cli framework gen-go orchestrator plugin-storage-markdown plugin-tools-features plugin-tools-marketplace plugin-transport-stdio proto sdk-go"

GREEN='\033[0;32m'; RED='\033[0;31m'; CYAN='\033[0;36m'; NC='\033[0m'
info()  { printf "${CYAN}[init]${NC}  %s\n" "$*"; }
ok()    { printf "${GREEN}[ok]${NC}    %s\n" "$*"; }
fail()  { printf "${RED}[fail]${NC}  %s\n" "$*" >&2; }

# Build list of new repos from orchestra.lock
REPOS="$(python3 -c "
import json
lock = json.load(open('${LOCK_FILE}'))
existing = set('${EXISTING}'.split())
for pkg in lock['packages']:
    name = pkg['name'].split('/')[-1]
    path = pkg['path']
    url = pkg['source']['url']
    if name not in existing:
        print(f'{name}|{path}|{url}')
")"

TOTAL=$(echo "$REPOS" | wc -l | tr -d ' ')
info "Initializing ${TOTAL} new repos..."
echo ""

SUCCESS=0
ERRORS=0

while IFS='|' read -r name path url; do
    [[ -z "$name" ]] && continue

    src="${ROOT}/${path}"
    tmp="$(mktemp -d)"

    info "${name}: initializing..."

    # Init a fresh git repo
    git init --quiet "${tmp}/repo"
    git -C "${tmp}/repo" config user.name "Fady Mondy"
    git -C "${tmp}/repo" config user.email "info@3x1.io"

    # Copy source files
    rsync -a --exclude '.git' "${src}/" "${tmp}/repo/"

    # Commit and push
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
done <<< "$REPOS"

echo ""
echo "────────────────────────────────"
ok "Initialized: ${SUCCESS}  Failed: ${ERRORS}"
[[ $ERRORS -gt 0 ]] && exit 1
exit 0
