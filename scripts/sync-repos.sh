#!/usr/bin/env bash
# sync-repos.sh — Push local libs/ changes to individual GitHub repos.
# Reads package list from orchestra.lock (single source of truth).
#
# Usage:
#   ./scripts/sync-repos.sh                    # Sync all repos
#   ./scripts/sync-repos.sh sdk-go cli         # Sync specific repos only
#   ./scripts/sync-repos.sh --dry-run          # Preview without pushing
#   ./scripts/sync-repos.sh --message "fix X"  # Custom commit message

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOCK_FILE="${ROOT}/orchestra.lock"
BRANCH="master"
DRY_RUN=false
MSG=""

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'
info()    { printf "${CYAN}[sync]${NC}  %s\n" "$*"; }
ok()      { printf "${GREEN}[ok]${NC}    %s\n" "$*"; }
warn()    { printf "${YELLOW}[skip]${NC}  %s\n" "$*"; }
fail()    { printf "${RED}[fail]${NC}  %s\n" "$*" >&2; }

if [[ ! -f "$LOCK_FILE" ]]; then
  fail "orchestra.lock not found at ${LOCK_FILE}"
  exit 1
fi

# Parse args
FILTER=()
while [[ $# -gt 0 ]]; do
  case "$1" in
    --dry-run)  DRY_RUN=true; shift ;;
    --message)  MSG="$2"; shift 2 ;;
    --help|-h)
      echo "Usage: $0 [--dry-run] [--message \"msg\"] [repo ...]"
      echo ""
      echo "Packages (from orchestra.lock):"
      python3 -c "
import json
lock = json.load(open('${LOCK_FILE}'))
for pkg in lock['packages']:
    name = pkg['name'].split('/')[-1]
    print(f'  {name:35s} {pkg[\"path\"]}')"
      exit 0
      ;;
    -*) fail "Unknown flag: $1"; exit 1 ;;
    *)  FILTER+=("$1"); shift ;;
  esac
done

# Read packages from orchestra.lock into temp file (one per line: name|path|url)
PKGLIST="$(mktemp)"
python3 -c "
import json
lock = json.load(open('${LOCK_FILE}'))
for pkg in lock['packages']:
    name = pkg['name'].split('/')[-1]
    path = pkg['path']
    url = pkg['source']['url']
    print(f'{name}|{path}|{url}')
" > "$PKGLIST"

# Build target list
TARGETLIST="$(mktemp)"
if [[ ${#FILTER[@]} -gt 0 ]]; then
  for filter in "${FILTER[@]}"; do
    match="$(grep "^${filter}|" "$PKGLIST" || true)"
    if [[ -z "$match" ]]; then
      fail "Unknown package: ${filter}"
      echo ""
      echo "Available packages:"
      cut -d'|' -f1 "$PKGLIST" | sed 's/^/  /'
      rm -f "$PKGLIST" "$TARGETLIST"
      exit 1
    fi
    echo "$match" >> "$TARGETLIST"
  done
else
  cp "$PKGLIST" "$TARGETLIST"
fi

rm -f "$PKGLIST"

# Preflight
command -v gh &>/dev/null || { fail "gh CLI not found"; rm -f "$TARGETLIST"; exit 1; }
command -v git &>/dev/null || { fail "git not found"; rm -f "$TARGETLIST"; exit 1; }

TIMESTAMP="$(date -u '+%Y-%m-%d %H:%M:%S UTC')"
COMMIT_MSG="${MSG:-Sync from monorepo (${TIMESTAMP})}"

count="$(wc -l < "$TARGETLIST" | tr -d ' ')"
info "Syncing ${count} package(s)"
$DRY_RUN && warn "DRY RUN — no changes will be pushed"
echo ""

SYNCED=0; UNCHANGED=0; ERRORS=0

while IFS='|' read -r name path url; do
  [[ -z "$name" ]] && continue

  src="${ROOT}/${path}"

  if [[ ! -d "$src" ]]; then
    fail "${name}: source dir not found (${path}/)"
    ERRORS=$((ERRORS + 1))
    continue
  fi

  # Clone into temp dir (full clone to have complete history)
  tmp="$(mktemp -d)"

  if ! git clone --quiet --branch "${BRANCH}" "${url}" "${tmp}/repo" 2>/dev/null; then
    fail "${name}: clone failed (${url})"
    rm -rf "$tmp"
    ERRORS=$((ERRORS + 1))
    continue
  fi

  cd "${tmp}/repo"

  # Pull latest to ensure we're on the tip of remote
  git pull --quiet origin "${BRANCH}" 2>/dev/null || true

  # Replace contents (keep .git)
  find "${tmp}/repo" -mindepth 1 -maxdepth 1 ! -name '.git' -exec rm -rf {} +
  cp -a "${src}/." "${tmp}/repo/"

  # Check for changes
  git add -A

  if git diff --cached --quiet; then
    warn "${name}: no changes"
    UNCHANGED=$((UNCHANGED + 1))
  else
    changes="$(git diff --cached --stat | tail -1)"
    info "${name}: ${changes}"

    if $DRY_RUN; then
      info "${name}: would push"
      SYNCED=$((SYNCED + 1))
    else
      git commit -m "${COMMIT_MSG}" --quiet
      if ! git push origin "${BRANCH}" --quiet 2>/dev/null; then
        # If push fails (race condition), pull --rebase and retry once
        warn "${name}: push failed, rebasing and retrying..."
        git pull --rebase --quiet origin "${BRANCH}" 2>/dev/null
        git push origin "${BRANCH}" --quiet
      fi
      ok "${name}: pushed"
      SYNCED=$((SYNCED + 1))
    fi
  fi

  cd "${ROOT}"
  rm -rf "$tmp"
done < "$TARGETLIST"

rm -f "$TARGETLIST"

# Summary
echo ""
echo "────────────────────────────────"
ok "Synced: ${SYNCED}  Unchanged: ${UNCHANGED}  Failed: ${ERRORS}"
[[ $ERRORS -gt 0 ]] && exit 1
exit 0
