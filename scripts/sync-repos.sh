#!/usr/bin/env bash
# sync-repos.sh — Push local libs/ changes to individual GitHub repos.
#
# Usage:
#   ./scripts/sync-repos.sh                    # Sync all repos
#   ./scripts/sync-repos.sh sdk-go cli         # Sync specific repos only
#   ./scripts/sync-repos.sh --dry-run          # Preview without pushing
#   ./scripts/sync-repos.sh --message "fix X"  # Custom commit message

set -euo pipefail

ORG="orchestra-mcp"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BRANCH="master"
DRY_RUN=false
MSG=""

# Dependency order — push base packages first
# Format: "name|source_dir"
ALL_REPOS=(
  "proto|libs/proto"
  "gen-go|libs/gen-go"
  "sdk-go|libs/sdk-go"
  "orchestrator|libs/orchestrator"
  "plugin-storage-markdown|libs/plugin-storage-markdown"
  "plugin-tools-features|libs/plugin-tools-features"
  "plugin-transport-stdio|libs/plugin-transport-stdio"
  "cli|libs/cli"
)

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'
info()    { printf "${CYAN}[sync]${NC}  %s\n" "$*"; }
ok()      { printf "${GREEN}[ok]${NC}    %s\n" "$*"; }
warn()    { printf "${YELLOW}[skip]${NC}  %s\n" "$*"; }
fail()    { printf "${RED}[fail]${NC}  %s\n" "$*" >&2; }

# Parse args
FILTER=()
while [[ $# -gt 0 ]]; do
  case "$1" in
    --dry-run)  DRY_RUN=true; shift ;;
    --message)  MSG="$2"; shift 2 ;;
    --help|-h)
      echo "Usage: $0 [--dry-run] [--message \"msg\"] [repo ...]"
      echo ""
      echo "Repos:"
      for entry in "${ALL_REPOS[@]}"; do
        IFS='|' read -r name _ <<< "$entry"
        echo "  $name"
      done
      exit 0
      ;;
    -*) fail "Unknown flag: $1"; exit 1 ;;
    *)  FILTER+=("$1"); shift ;;
  esac
done

# Build target list
TARGETS=()
if [[ ${#FILTER[@]} -gt 0 ]]; then
  for filter in "${FILTER[@]}"; do
    found=false
    for entry in "${ALL_REPOS[@]}"; do
      IFS='|' read -r name srcdir <<< "$entry"
      if [[ "$name" == "$filter" ]]; then
        TARGETS+=("$entry")
        found=true
        break
      fi
    done
    if ! $found; then
      fail "Unknown repo: $filter"
      exit 1
    fi
  done
else
  TARGETS=("${ALL_REPOS[@]}")
fi

# Preflight
command -v gh &>/dev/null || { fail "gh CLI not found"; exit 1; }
command -v git &>/dev/null || { fail "git not found"; exit 1; }

TIMESTAMP="$(date -u '+%Y-%m-%d %H:%M:%S UTC')"
COMMIT_MSG="${MSG:-Sync from monorepo (${TIMESTAMP})}"

info "Syncing ${#TARGETS[@]} repo(s) to ${ORG}/"
$DRY_RUN && warn "DRY RUN — no changes will be pushed"
echo ""

SYNCED=0; UNCHANGED=0; ERRORS=0

for entry in "${TARGETS[@]}"; do
  IFS='|' read -r name srcdir <<< "$entry"
  src="${ROOT}/${srcdir}"

  if [[ ! -d "$src" ]]; then
    fail "${name}: source dir not found (${srcdir}/)"
    ERRORS=$((ERRORS + 1))
    continue
  fi

  # Clone into temp dir
  tmp="$(mktemp -d)"

  if ! git clone --quiet --branch "${BRANCH}" "https://github.com/${ORG}/${name}.git" "${tmp}/repo" 2>/dev/null; then
    fail "${name}: clone failed"
    rm -rf "$tmp"
    ERRORS=$((ERRORS + 1))
    continue
  fi

  # Replace contents (keep .git)
  find "${tmp}/repo" -mindepth 1 -maxdepth 1 ! -name '.git' -exec rm -rf {} +
  cp -a "${src}/." "${tmp}/repo/"

  # Check for changes
  cd "${tmp}/repo"
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
      git push origin "${BRANCH}" --quiet
      ok "${name}: pushed"
      SYNCED=$((SYNCED + 1))
    fi
  fi

  cd "${ROOT}"
  rm -rf "$tmp"
done

# Summary
echo ""
echo "────────────────────────────────"
ok "Synced: ${SYNCED}  Unchanged: ${UNCHANGED}  Failed: ${ERRORS}"
[[ $ERRORS -gt 0 ]] && exit 1
exit 0
