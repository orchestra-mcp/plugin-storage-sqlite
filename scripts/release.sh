#!/usr/bin/env bash
# release.sh — Sync libs/ to GitHub repos, then tag + create GitHub releases.
# Reads package list from orchestra.lock (single source of truth).
#
# Usage:
#   ./scripts/release.sh v0.2.0                          # Release all repos at v0.2.0
#   ./scripts/release.sh v0.2.0 sdk-go cli               # Release specific repos only
#   ./scripts/release.sh v0.2.0 --dry-run                # Preview without pushing
#   ./scripts/release.sh v0.2.0 --force                  # Delete existing tag and re-release
#   ./scripts/release.sh v0.2.0 --sync-only              # Sync code but skip tag/release
#   ./scripts/release.sh v0.2.0 --skip-sync              # Tag/release only (code already synced)

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOCK_FILE="${ROOT}/orchestra.lock"
BRANCH="master"
DRY_RUN=false
FORCE=false
SYNC_ONLY=false
SKIP_SYNC=false

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'
info()    { printf "${CYAN}[release]${NC} %s\n" "$*"; }
ok()      { printf "${GREEN}[ok]${NC}      %s\n" "$*"; }
warn()    { printf "${YELLOW}[skip]${NC}    %s\n" "$*"; }
fail()    { printf "${RED}[fail]${NC}    %s\n" "$*" >&2; }

if [[ ! -f "$LOCK_FILE" ]]; then
  fail "orchestra.lock not found at ${LOCK_FILE}"
  exit 1
fi

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
      echo "Packages (from orchestra.lock):"
      python3 -c "
import json
lock = json.load(open('${LOCK_FILE}'))
for pkg in lock['packages']:
    name = pkg['name'].split('/')[-1]
    print(f'  {name:35s} {pkg[\"description\"][:60]}')"
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

# Read packages from orchestra.lock (one per line: name|url|description)
PKGLIST="$(mktemp)"
python3 -c "
import json
lock = json.load(open('${LOCK_FILE}'))
org = lock['packages'][0]['name'].split('/')[0]
for pkg in lock['packages']:
    name = pkg['name'].split('/')[-1]
    url = pkg['source']['url'].replace('.git', '')
    desc = pkg['description']
    print(f'{name}|{org}/{name}|{desc}')
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

count="$(wc -l < "$TARGETLIST" | tr -d ' ')"

echo ""
info "Release ${VERSION} for ${count} package(s)"
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
  while IFS='|' read -r name _ _; do
    SYNC_NAMES+=("$name")
  done < "$TARGETLIST"

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
  rm -f "$TARGETLIST"
  ok "Sync complete (--sync-only, skipping tags)"
  exit 0
fi

# ──────────────────────────────────────────────────────
# Step 2: Tag and create GitHub releases
# ──────────────────────────────────────────────────────

info "Step 2: Creating tags and releases..."
echo ""

RELEASED=0; SKIPPED=0; ERRORS=0

while IFS='|' read -r name repo desc; do
  [[ -z "$name" ]] && continue

  # Check if release already exists
  if gh release view "${VERSION}" --repo "${repo}" &>/dev/null 2>&1; then
    if $FORCE; then
      info "${name}: deleting existing ${VERSION}"
      if ! $DRY_RUN; then
        gh release delete "${VERSION}" --repo "${repo}" --yes 2>/dev/null || true
        gh api -X DELETE "repos/${repo}/git/refs/tags/${VERSION}" 2>/dev/null || true
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
    --notes "${desc}" \
    --target "${BRANCH}" 2>/dev/null; then
    ok "${name}: https://github.com/${repo}/releases/tag/${VERSION}"
    RELEASED=$((RELEASED + 1))
  else
    fail "${name}: release creation failed"
    ERRORS=$((ERRORS + 1))
  fi
done < "$TARGETLIST"

rm -f "$TARGETLIST"

# ──────────────────────────────────────────────────────
# Step 3: Update orchestra.lock with new version
# ──────────────────────────────────────────────────────

if [[ $RELEASED -gt 0 ]] && ! $DRY_RUN; then
  info "Step 3: Updating orchestra.lock to ${VERSION}..."
  python3 -c "
import json
lock = json.load(open('${LOCK_FILE}'))
for pkg in lock['packages']:
    pkg['version'] = '${VERSION}'
with open('${LOCK_FILE}', 'w') as f:
    json.dump(lock, f, indent=4)
    f.write('\n')
"
  ok "orchestra.lock updated to ${VERSION}"
fi

# ──────────────────────────────────────────────────────
# Summary
# ──────────────────────────────────────────────────────

echo ""
echo "────────────────────────────────"
ok "Released: ${RELEASED}  Skipped: ${SKIPPED}  Failed: ${ERRORS}"

if [[ $RELEASED -gt 0 ]] && ! $DRY_RUN; then
  org="$(python3 -c "import json; print(json.load(open('${LOCK_FILE}'))['packages'][0]['name'].split('/')[0])")"
  echo ""
  info "Go module proxy will index tags within a few minutes."
  info "Verify: GOWORK=off go list -m github.com/${org}/sdk-go@${VERSION}"
fi

[[ $ERRORS -gt 0 ]] && exit 1
exit 0
