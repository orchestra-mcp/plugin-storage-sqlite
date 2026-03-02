#!/bin/sh
#
# Orchestra MCP — Install Script
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/orchestra-mcp/framework/master/scripts/install.sh | sh
#
# Environment variables:
#   INSTALL_DIR   — Installation directory (default: /usr/local/bin)
#   VERSION       — Version to install (default: latest)
#   GITHUB_REPO   — GitHub repository (default: orchestra-mcp/framework)
#
set -e

GITHUB_REPO="${GITHUB_REPO:-orchestra-mcp/framework}"
VERSION="${VERSION:-latest}"
EXE=""
CURL_EXTRA=""

# Detect OS and architecture.
detect_platform() {
    RAW_OS="$(uname -s)"
    ARCH="$(uname -m)"

    case "$RAW_OS" in
        Darwin*)          OS="darwin" ;;
        Linux*)           OS="linux" ;;
        MINGW*|MSYS*|CYGWIN*)
            OS="windows"
            EXE=".exe"
            # Windows Schannel may fail revocation checks behind proxies/firewalls.
            CURL_EXTRA="--ssl-revoke-best-effort"
            ;;
        *)
            echo "Error: Unsupported OS: $RAW_OS" >&2
            echo "Orchestra supports macOS, Linux, and Windows (Git Bash)." >&2
            exit 1
            ;;
    esac

    case "$ARCH" in
        x86_64|amd64)  ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        *)
            echo "Error: Unsupported architecture: $ARCH" >&2
            echo "Orchestra supports amd64 and arm64." >&2
            exit 1
            ;;
    esac

    # Default install dir per platform.
    if [ -z "${INSTALL_DIR:-}" ]; then
        if [ "$OS" = "windows" ]; then
            INSTALL_DIR="${LOCALAPPDATA:-$HOME/AppData/Local}/Orchestra/bin"
        else
            INSTALL_DIR="/usr/local/bin"
        fi
    fi

    PLATFORM="${OS}-${ARCH}"
}

# Resolve the download URL.
resolve_url() {
    if [ "$VERSION" = "latest" ]; then
        # GitHub /releases/latest only returns non-prerelease. Use API to get the
        # most recent release (including prereleases).
        if command -v curl >/dev/null 2>&1; then
            VERSION="$(curl -fsSL $CURL_EXTRA "https://api.github.com/repos/${GITHUB_REPO}/releases" \
                | grep -m1 '"tag_name"' | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')"
        elif command -v wget >/dev/null 2>&1; then
            VERSION="$(wget -qO- "https://api.github.com/repos/${GITHUB_REPO}/releases" \
                | grep -m1 '"tag_name"' | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')"
        fi
        if [ -z "$VERSION" ]; then
            echo "Error: Could not determine latest version." >&2
            exit 1
        fi
    fi
    URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/orchestra-${PLATFORM}.tar.gz"
}

# Download and install.
install() {
    echo "Orchestra MCP Installer"
    echo "======================"
    echo ""
    echo "  Platform:    ${PLATFORM}"
    echo "  Install dir: ${INSTALL_DIR}"
    echo "  Version:     ${VERSION}"
    echo ""

    # Create temp directory.
    TMP_DIR="$(mktemp -d)"
    trap 'rm -rf "$TMP_DIR"' EXIT

    echo "Downloading ${URL}..."
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL $CURL_EXTRA "$URL" -o "${TMP_DIR}/orchestra.tar.gz"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$URL" -O "${TMP_DIR}/orchestra.tar.gz"
    else
        echo "Error: curl or wget is required." >&2
        exit 1
    fi

    echo "Extracting..."
    tar -xzf "${TMP_DIR}/orchestra.tar.gz" -C "$TMP_DIR"

    echo "Installing to ${INSTALL_DIR}..."

    # Check write permissions, use sudo if needed (skip sudo on Windows).
    SUDO=""
    if [ "$OS" != "windows" ]; then
        if [ ! -w "$INSTALL_DIR" ] 2>/dev/null || { [ ! -d "$INSTALL_DIR" ] && ! mkdir -p "$INSTALL_DIR" 2>/dev/null; }; then
            if command -v sudo >/dev/null 2>&1; then
                echo "  (need sudo for ${INSTALL_DIR})"
                SUDO="sudo"
            else
                echo "Error: No write permission to ${INSTALL_DIR} and sudo not available." >&2
                echo "Try: INSTALL_DIR=~/.local/bin sh install.sh" >&2
                exit 1
            fi
        fi
    fi

    $SUDO mkdir -p "$INSTALL_DIR"

    for bin in orchestra orchestrator storage-markdown tools-features transport-stdio tools-marketplace; do
        if [ -f "${TMP_DIR}/${bin}${EXE}" ]; then
            $SUDO cp "${TMP_DIR}/${bin}${EXE}" "${INSTALL_DIR}/${bin}${EXE}"
            if [ "$OS" != "windows" ]; then
                $SUDO chmod +x "${INSTALL_DIR}/${bin}${EXE}"
            fi
            echo "  installed ${INSTALL_DIR}/${bin}${EXE}"
        fi
    done

    echo ""
    echo "Done! Orchestra MCP installed."
    echo ""
    echo "Next steps:"
    echo "  cd your-project"
    echo "  orchestra init"
    echo ""

    # Verify it works.
    if command -v "orchestra${EXE}" >/dev/null 2>&1; then
        echo "Version: $(orchestra${EXE} version)"
    else
        echo "Note: Add ${INSTALL_DIR} to your PATH if not already there."
        if [ "$OS" = "windows" ]; then
            echo ""
            echo "  Run this in PowerShell (as Administrator):"
            echo "    [Environment]::SetEnvironmentVariable('Path', [Environment]::GetEnvironmentVariable('Path', 'User') + ';${INSTALL_DIR}', 'User')"
            echo ""
            echo "  Then restart your terminal."
        fi
    fi
}

detect_platform
resolve_url
install
