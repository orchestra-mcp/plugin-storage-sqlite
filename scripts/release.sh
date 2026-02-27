#!/usr/bin/env bash
# release.sh — Sync libs/ to GitHub repos, then tag + create GitHub releases.
#
# Usage:
#   ./scripts/release.sh v0.2.0                          # Release all repos at v0.2.0
#   ./scripts/release.sh v0.2.0 sdk-go cli               # Release specific repos only
#   ./scripts/release.sh v0.2.0 --dry-run                # Preview without pushing
#   ./scripts/release.sh v0.2.0 --force                  # Delete existing tag and re-release
#   ./scripts/release.sh v0.2.0 --sync-only              # Sync code but skip tag/release
#   ./scripts/release.sh v0.2.0 --skip-sync              # Tag/release only (code already synced)

set -euo pipefail

ORG="orchestra-mcp"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BRANCH="master"
DRY_RUN=false
FORCE=false
SYNC_ONLY=false
SKIP_SYNC=false

# Dependency order — release base packages first so dependents can resolve them
# Format: "name|release_note"
ALL_REPOS=(
  "proto|Protobuf contract definitions for the Orchestra MCP plugin host protocol."
  "gen-go|Generated Go code from Orchestra MCP Protobuf definitions."
  "sdk-go|Go SDK for building Orchestra MCP plugins (QUIC transport, mTLS, Protobuf framing)."
  "orchestrator|Plugin host orchestrator with QUIC mesh, message routing, and lifecycle management."
  "plugin-storage-markdown|Markdown storage plugin with YAML frontmatter metadata."
  "plugin-tools-features|Feature-driven project management tools (workflow, dependencies, WIP limits, reviews)."
  "plugin-transport-stdio|MCP stdio transport plugin (JSON-RPC over stdin/stdout)."
  "cli|Orchestra CLI for installing and managing plugins."
)

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'
info()    { printf "${CYAN}[release]${NC} %s\n" "$*"; }
ok()      { printf "${GREEN}[ok]${NC}      %s\n" "$*"; }
warn()    { printf "${YELLOW}[skip]${NC}    %s\n" "$*"; }
fail()    { printf "${RED}[fail]${NC}    %s\n" "$*" >&2; }

# Parse args
VERSION=""
FILTER=()
while [[ $# -gt 0 ]]; do
  case "$1" in
    --dry-run)    DRY_RUN=true; shift ;;
    --force)      FORCE=true; shift ;;
    --sync-only)  SYNC_ONLY=true; shift ;;
    --skip-sync)  SKIP_SYNC=true; shift ;;
    --help|-h)
      echo "Usage: $0 <version> [--dry-run] [--force] [--sync-only] [--skip-sync] [repo ...]"
      echo ""
      echo "  version       Semver tag (e.g. v0.2.0)"
      echo "  --dry-run     Preview without pushing"
      echo "  --force       Delete existing tags and re-release"
      echo "  --sync-only   Sync code only, skip tagging"
      echo "  --skip-sync   Skip sync, tag/release only"
      echo ""
      echo "Repos:"
      for entry in "${ALL_REPOS[@]}"; do
        IFS='|' read -r name _ <<< "$entry"
        echo "  $name"
      done
      exit 0
      ;;
    v*)  VERSION="$1"; shift ;;
    -*)  fail "Unknown flag: $1"; exit 1 ;;
    *)   FILTER+=("$1"); shift ;;
  esac
done

if [[ -z "$VERSION" ]]; then
  fail "Version required. Usage: $0 v0.2.0 [repos...]"
  exit 1
fi

# Validate semver format
if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  fail "Invalid version format: $VERSION (expected vX.Y.Z)"
  exit 1
fi

# Build target list (names only for filtering, full entries for release notes)
TARGETS=()
if [[ ${#FILTER[@]} -gt 0 ]]; then
  for filter in "${FILTER[@]}"; do
    found=false
    for entry in "${ALL_REPOS[@]}"; do
      IFS='|' read -r name note <<< "$entry"
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

echo ""
info "Release ${VERSION} for ${#TARGETS[@]} repo(s)"
$DRY_RUN && warn "DRY RUN — no changes will be pushed"
echo ""

# ──────────────────────────────────────────────────────
# Step 1: Sync code to GitHub repos
# ──────────────────────────────────────────────────────

if ! $SKIP_SYNC; then
  info "Step 1: Syncing code to GitHub..."
  echo ""

  # Extract just the names for sync script
  SYNC_NAMES=()
  for entry in "${TARGETS[@]}"; do
    IFS='|' read -r name _ <<< "$entry"
    SYNC_NAMES+=("$name")
  done

  SYNC_ARGS=()
  $DRY_RUN && SYNC_ARGS+=(--dry-run)
  SYNC_ARGS+=(--message "Release ${VERSION}")
  SYNC_ARGS+=("${SYNC_NAMES[@]}")

  "${ROOT}/scripts/sync-repos.sh" "${SYNC_ARGS[@]}"
  echo ""
else
  warn "Skipping sync (--skip-sync)"
  echo ""
fi

if $SYNC_ONLY; then
  ok "Sync complete (--sync-only, skipping tags)"
  exit 0
fi

# ──────────────────────────────────────────────────────
# Step 2: Tag and create GitHub releases
# ──────────────────────────────────────────────────────

info "Step 2: Creating tags and releases..."
echo ""

RELEASED=0; SKIPPED=0; ERRORS=0

for entry in "${TARGETS[@]}"; do
  IFS='|' read -r name note <<< "$entry"
  repo="${ORG}/${name}"

  # Check if release already exists
  if gh release view "${VERSION}" --repo "${repo}" &>/dev/null 2>&1; then
    if $FORCE; then
      info "${name}: deleting existing ${VERSION}"
      if ! $DRY_RUN; then
        gh release delete "${VERSION}" --repo "${repo}" --yes 2>/dev/null || true
        gh api -X DELETE "repos/${repo}/git/refs/tags/${VERSION}" 2>/dev/null || true
        # Small delay to let GitHub process the deletion
        sleep 1
      fi
    else
      warn "${name}: ${VERSION} already exists (use --force to overwrite)"
      SKIPPED=$((SKIPPED + 1))
      continue
    fi
  fi

  if $DRY_RUN; then
    info "${name}: would create release ${VERSION}"
    RELEASED=$((RELEASED + 1))
    continue
  fi

  # Create the release (this also creates the tag)
  if gh release create "${VERSION}" \
    --repo "${repo}" \
    --title "${VERSION}" \
    --notes "${note}" \
    --target "${BRANCH}" 2>/dev/null; then
    ok "${name}: https://github.com/${repo}/releases/tag/${VERSION}"
    RELEASED=$((RELEASED + 1))
  else
    fail "${name}: release creation failed"
    ERRORS=$((ERRORS + 1))
  fi
done

# ──────────────────────────────────────────────────────
# Summary
# ──────────────────────────────────────────────────────

echo ""
echo "────────────────────────────────"
ok "Released: ${RELEASED}  Skipped: ${SKIPPED}  Failed: ${ERRORS}"

if [[ $RELEASED -gt 0 ]] && ! $DRY_RUN; then
  echo ""
  info "Go module proxy will index tags within a few minutes."
  info "Verify: GOWORK=off go list -m github.com/${ORG}/sdk-go@${VERSION}"
fi

[[ $ERRORS -gt 0 ]] && exit 1
exit 0
