#!/usr/bin/env bash
# split-repos.sh — Create GitHub repos and push initial code from monorepo directories.
#
# Usage:
#   ./scripts/split-repos.sh              # Create repos and push code
#   ./scripts/split-repos.sh --dry-run    # Show what would be created without doing it
#
# Prerequisites:
#   - gh CLI installed and authenticated
#   - git installed
#   - Current directory is the monorepo root

set -euo pipefail

# ---------------------------------------------------------------------------
# Configuration
# ---------------------------------------------------------------------------

ORG="orchestra-mcp"
MONOREPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DRY_RUN=false
BRANCH="master"

# Each entry: "repo_name|source_dir|description"
# source_dir is relative to MONOREPO_ROOT
REPOS=(
  "proto|libs/proto|Protobuf schema for Orchestra plugin protocol"
  "gen-go|libs/gen-go|Generated Go protobuf code"
  "sdk-go|libs/sdk-go|Plugin SDK for Go"
  "orchestrator|libs/orchestrator|Central hub service"
  "plugin-storage-markdown|libs/plugin-storage-markdown|Markdown file storage plugin"
  "plugin-tools-features|libs/plugin-tools-features|34 workflow tools plugin"
  "plugin-transport-stdio|libs/plugin-transport-stdio|MCP JSON-RPC bridge plugin"
  "cli|libs/cli|Orchestra CLI binary"
)

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

info()    { printf "${CYAN}[INFO]${NC}  %s\n" "$*"; }
success() { printf "${GREEN}[OK]${NC}    %s\n" "$*"; }
warn()    { printf "${YELLOW}[WARN]${NC}  %s\n" "$*"; }
error()   { printf "${RED}[ERROR]${NC} %s\n" "$*" >&2; }

cleanup() {
  if [[ -n "${TMPDIR_WORK:-}" && -d "${TMPDIR_WORK}" ]]; then
    rm -rf "${TMPDIR_WORK}"
  fi
}
trap cleanup EXIT

# ---------------------------------------------------------------------------
# Parse arguments
# ---------------------------------------------------------------------------

for arg in "$@"; do
  case "$arg" in
    --dry-run) DRY_RUN=true ;;
    --help|-h)
      echo "Usage: $0 [--dry-run]"
      echo ""
      echo "Creates GitHub repos under ${ORG}/ and pushes initial code."
      echo ""
      echo "Options:"
      echo "  --dry-run   Show what would be created without doing it"
      echo "  --help      Show this help message"
      exit 0
      ;;
    *)
      error "Unknown argument: $arg"
      exit 1
      ;;
  esac
done

# ---------------------------------------------------------------------------
# Preflight checks
# ---------------------------------------------------------------------------

info "Monorepo root: ${MONOREPO_ROOT}"

if $DRY_RUN; then
  warn "DRY RUN MODE — no repos will be created or pushed"
  echo ""
fi

# Check gh CLI
if ! command -v gh &>/dev/null; then
  error "gh CLI not found. Install it: https://cli.github.com/"
  exit 1
fi

if ! gh auth status &>/dev/null; then
  error "gh CLI not authenticated. Run: gh auth login"
  exit 1
fi
success "gh CLI authenticated"

# Check git
if ! command -v git &>/dev/null; then
  error "git not found"
  exit 1
fi

# Verify source directories exist
for entry in "${REPOS[@]}"; do
  IFS='|' read -r name srcdir desc <<< "$entry"
  if [[ ! -d "${MONOREPO_ROOT}/${srcdir}" ]]; then
    error "Source directory not found: ${MONOREPO_ROOT}/${srcdir}"
    exit 1
  fi
done
success "All source directories verified"

echo ""

# ---------------------------------------------------------------------------
# Create and push repos
# ---------------------------------------------------------------------------

CREATED=()
SKIPPED=()
FAILED=()

for entry in "${REPOS[@]}"; do
  IFS='|' read -r name srcdir desc <<< "$entry"

  echo "──────────────────────────────────────────────────────────"
  info "Repository: ${ORG}/${name}"
  info "  Source:      ${srcdir}/"
  info "  Description: ${desc}"

  if $DRY_RUN; then
    info "  [dry-run] Would create ${ORG}/${name} and push ${srcdir}/"
    CREATED+=("$name")
    echo ""
    continue
  fi

  # Check if repo already exists
  if gh repo view "${ORG}/${name}" &>/dev/null; then
    warn "  Repo ${ORG}/${name} already exists — skipping creation"
    SKIPPED+=("$name")
    echo ""
    continue
  fi

  # Create the GitHub repo
  info "  Creating GitHub repo..."
  if ! gh repo create "${ORG}/${name}" --public --description "${desc}" 2>&1; then
    error "  Failed to create repo ${ORG}/${name}"
    FAILED+=("$name")
    echo ""
    continue
  fi
  success "  Repo created: ${ORG}/${name}"

  # Prepare temp directory
  TMPDIR_WORK="$(mktemp -d)"
  WORK="${TMPDIR_WORK}/${name}"
  mkdir -p "${WORK}"

  # Copy source directory contents
  info "  Copying ${srcdir}/ to temp directory..."
  cp -a "${MONOREPO_ROOT}/${srcdir}/." "${WORK}/"

  # Initialize git and push
  info "  Initializing git and pushing..."
  (
    cd "${WORK}"
    git init -b "${BRANCH}" --quiet
    git add -A
    git commit -m "Initial commit from monorepo split" --quiet
    git remote add origin "https://github.com/${ORG}/${name}.git"
    git push -u origin "${BRANCH}" --quiet
  )

  if [[ $? -eq 0 ]]; then
    success "  Pushed to ${ORG}/${name}"
    CREATED+=("$name")
  else
    error "  Failed to push ${ORG}/${name}"
    FAILED+=("$name")
  fi

  # Clean up temp dir for this iteration
  rm -rf "${TMPDIR_WORK}"
  unset TMPDIR_WORK

  echo ""
done

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------

echo "══════════════════════════════════════════════════════════"
echo ""

if $DRY_RUN; then
  info "DRY RUN SUMMARY"
  echo ""
  info "Would create ${#CREATED[@]} repos under ${ORG}/:"
  for name in "${CREATED[@]}"; do
    echo "  - https://github.com/${ORG}/${name}"
  done
  echo ""
  info "Meta-repo (orchestra) is not included — set up separately as the root."
  echo ""
  exit 0
fi

if [[ ${#CREATED[@]} -gt 0 ]]; then
  success "Created ${#CREATED[@]} repo(s):"
  for name in "${CREATED[@]}"; do
    echo "  - https://github.com/${ORG}/${name}"
  done
fi

if [[ ${#SKIPPED[@]} -gt 0 ]]; then
  warn "Skipped ${#SKIPPED[@]} repo(s) (already exist):"
  for name in "${SKIPPED[@]}"; do
    echo "  - https://github.com/${ORG}/${name}"
  done
fi

if [[ ${#FAILED[@]} -gt 0 ]]; then
  error "Failed ${#FAILED[@]} repo(s):"
  for name in "${FAILED[@]}"; do
    echo "  - ${name}"
  done
  echo ""
  exit 1
fi

echo ""
info "Meta-repo (orchestra) is not included — set up separately as the root."
info "Next step: run ./scripts/sync-repos.sh to keep repos in sync."
echo ""
