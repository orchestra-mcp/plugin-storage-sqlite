#!/usr/bin/env bash
# orchestra-fmt.sh — Format and validate orchestra.json and orchestra.lock files.
#
# Usage:
#   ./scripts/orchestra-fmt.sh              # Format all orchestra.json + orchestra.lock
#   ./scripts/orchestra-fmt.sh --check      # Check only (exit 1 if changes needed)
#   ./scripts/orchestra-fmt.sh --validate   # Validate structure + cross-references

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CHECK=false
VALIDATE=false

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'
info()  { printf "${CYAN}[fmt]${NC}   %s\n" "$*"; }
ok()    { printf "${GREEN}[ok]${NC}    %s\n" "$*"; }
warn()  { printf "${YELLOW}[warn]${NC}  %s\n" "$*"; }
fail()  { printf "${RED}[fail]${NC}  %s\n" "$*" >&2; }

while [[ $# -gt 0 ]]; do
  case "$1" in
    --check)    CHECK=true; shift ;;
    --validate) VALIDATE=true; shift ;;
    --help|-h)
      echo "Usage: $0 [--check] [--validate]"
      echo ""
      echo "  --check      Check formatting only (exit 1 if changes needed)"
      echo "  --validate   Validate structure and cross-references"
      echo ""
      echo "Without flags: format all orchestra.json and orchestra.lock files."
      exit 0
      ;;
    *) fail "Unknown argument: $1"; exit 1 ;;
  esac
done

ERRORS=0
FORMATTED=0

# ──────────────────────────────────────────────────────
# Format a JSON file (sort keys consistently, 4-space indent)
# ──────────────────────────────────────────────────────
format_json() {
  local file="$1"
  local label="$2"

  if [[ ! -f "$file" ]]; then
    return
  fi

  # Format with sorted keys and 4-space indent
  formatted="$(python3 -c "
import json, sys
with open('${file}') as f:
    data = json.load(f)
print(json.dumps(data, indent=4, ensure_ascii=False))
" 2>/dev/null)" || {
    fail "${label}: invalid JSON"
    ERRORS=$((ERRORS + 1))
    return
  }

  current="$(cat "$file")"

  # Compare (strip trailing newline for comparison)
  if [[ "$formatted" != "${current%$'\n'}" ]]; then
    if $CHECK; then
      fail "${label}: needs formatting"
      ERRORS=$((ERRORS + 1))
    else
      echo "$formatted" > "$file"
      ok "${label}: formatted"
      FORMATTED=$((FORMATTED + 1))
    fi
  else
    ok "${label}: ok"
  fi
}

# ──────────────────────────────────────────────────────
# Format root files
# ──────────────────────────────────────────────────────

info "Checking root files..."
format_json "${ROOT}/orchestra.json" "orchestra.json"
format_json "${ROOT}/orchestra.lock" "orchestra.lock"

# ──────────────────────────────────────────────────────
# Format plugin orchestra.json files
# ──────────────────────────────────────────────────────

info "Checking plugin files..."
for dir in "${ROOT}"/libs/*/; do
  pfile="${dir}orchestra.json"
  if [[ -f "$pfile" ]]; then
    name="$(basename "$dir")"
    format_json "$pfile" "libs/${name}/orchestra.json"
  fi
done

# ──────────────────────────────────────────────────────
# Validate (if requested)
# ──────────────────────────────────────────────────────

if $VALIDATE; then
  echo ""
  info "Validating structure and cross-references..."

  python3 -c "
import json, sys, os

root = '${ROOT}'
errors = []

# Load root manifest
try:
    with open(f'{root}/orchestra.json') as f:
        manifest = json.load(f)
except Exception as e:
    print(f'  FAIL  orchestra.json: {e}')
    sys.exit(1)

# Load lock file
try:
    with open(f'{root}/orchestra.lock') as f:
        lock = json.load(f)
except Exception as e:
    print(f'  FAIL  orchestra.lock: {e}')
    sys.exit(1)

# Validate required fields in manifest
for field in ['name', 'description', 'type', 'require', 'config']:
    if field not in manifest:
        errors.append(f'orchestra.json: missing required field \"{field}\"')

# Validate lock structure
if 'packages' not in lock:
    errors.append('orchestra.lock: missing \"packages\" array')
else:
    lock_names = set()
    for i, pkg in enumerate(lock['packages']):
        for field in ['name', 'version', 'source', 'type', 'description', 'path']:
            if field not in pkg:
                errors.append(f'orchestra.lock: packages[{i}] missing \"{field}\"')
        if 'name' in pkg:
            lock_names.add(pkg['name'])
        if 'source' in pkg:
            src = pkg['source']
            for sf in ['type', 'url', 'reference']:
                if sf not in src:
                    errors.append(f'orchestra.lock: packages[{i}].source missing \"{sf}\"')

    # Cross-reference: every require in manifest must be in lock
    for req_name in manifest.get('require', {}):
        if req_name not in lock_names:
            errors.append(f'orchestra.json requires \"{req_name}\" but it is not in orchestra.lock')

    # Cross-reference: every lock package should have matching plugin orchestra.json
    for pkg in lock['packages']:
        path = pkg.get('path', '')
        pfile = f'{root}/{path}/orchestra.json'
        if os.path.isdir(f'{root}/{path}') and not os.path.isfile(pfile):
            errors.append(f'{path}/orchestra.json: missing (expected for lock package \"{pkg[\"name\"]}\")')
        elif os.path.isfile(pfile):
            try:
                with open(pfile) as f:
                    pdata = json.load(f)
                if pdata.get('name') != pkg['name']:
                    errors.append(f'{path}/orchestra.json: name \"{pdata.get(\"name\")}\" does not match lock \"{pkg[\"name\"]}\"')
            except Exception as e:
                errors.append(f'{path}/orchestra.json: {e}')

    # Cross-reference: install-order matches lock packages
    order = manifest.get('extra', {}).get('install-order', [])
    if order:
        for item in order:
            full_name = f'{manifest.get(\"config\", {}).get(\"org\", \"orchestra-mcp\")}/{item}'
            if full_name not in lock_names:
                errors.append(f'orchestra.json install-order contains \"{item}\" but \"{full_name}\" not in lock')

# Plugin-level validation
for pkg in lock.get('packages', []):
    path = pkg.get('path', '')
    pfile = f'{root}/{path}/orchestra.json'
    if os.path.isfile(pfile):
        with open(pfile) as f:
            pdata = json.load(f)
        for field in ['name', 'description', 'type', 'license']:
            if field not in pdata:
                errors.append(f'{path}/orchestra.json: missing required field \"{field}\"')

if errors:
    for e in errors:
        print(f'  FAIL  {e}')
    sys.exit(1)
else:
    print('  All validations passed.')
    sys.exit(0)
" 2>&1 | while IFS= read -r line; do
    echo "$line"
  done

  if [[ ${PIPESTATUS[0]} -ne 0 ]]; then
    ERRORS=$((ERRORS + 1))
  else
    ok "Validation passed"
  fi
fi

# ──────────────────────────────────────────────────────
# Summary
# ──────────────────────────────────────────────────────

echo ""
if [[ $ERRORS -gt 0 ]]; then
  if $CHECK; then
    fail "Formatting check failed (${ERRORS} file(s) need formatting)"
  else
    fail "${ERRORS} error(s) found"
  fi
  exit 1
else
  if $CHECK; then
    ok "All files properly formatted"
  elif [[ $FORMATTED -gt 0 ]]; then
    ok "Formatted ${FORMATTED} file(s)"
  else
    ok "All files already formatted"
  fi
fi
