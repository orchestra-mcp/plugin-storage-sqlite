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

INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
GITHUB_REPO="${GITHUB_REPO:-orchestra-mcp/framework}"
VERSION="${VERSION:-latest}"

# Detect OS and architecture.
detect_platform() {
    OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
    ARCH="$(uname -m)"

    case "$OS" in
        darwin) OS="darwin" ;;
        linux)  OS="linux" ;;
        *)
            echo "Error: Unsupported OS: $OS" >&2
            echo "Orchestra supports macOS and Linux." >&2
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

    PLATFORM="${OS}-${ARCH}"
}

# Resolve the download URL.
resolve_url() {
    if [ "$VERSION" = "latest" ]; then
        URL="https://github.com/${GITHUB_REPO}/releases/latest/download/orchestra-${PLATFORM}.tar.gz"
    else
        URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/orchestra-${PLATFORM}.tar.gz"
    fi
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
        curl -fsSL "$URL" -o "${TMP_DIR}/orchestra.tar.gz"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$URL" -O "${TMP_DIR}/orchestra.tar.gz"
    else
        echo "Error: curl or wget is required." >&2
        exit 1
    fi

    echo "Extracting..."
    tar -xzf "${TMP_DIR}/orchestra.tar.gz" -C "$TMP_DIR"

    echo "Installing to ${INSTALL_DIR}..."
    mkdir -p "$INSTALL_DIR"

    for bin in orchestra orchestrator storage-markdown tools-features transport-stdio tools-marketplace; do
        if [ -f "${TMP_DIR}/${bin}" ]; then
            cp "${TMP_DIR}/${bin}" "${INSTALL_DIR}/${bin}"
            chmod +x "${INSTALL_DIR}/${bin}"
            echo "  installed ${INSTALL_DIR}/${bin}"
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
    if command -v orchestra >/dev/null 2>&1; then
        echo "Version: $(orchestra version)"
    else
        echo "Note: Add ${INSTALL_DIR} to your PATH if not already there."
    fi
}

detect_platform
resolve_url
install
