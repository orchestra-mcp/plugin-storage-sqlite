# Homebrew formula for Orchestra MCP
#
# To use this formula, create a tap repository:
#   https://github.com/orchestra-mcp/homebrew-tap
#
# Then users can install with:
#   brew tap orchestra-mcp/tap
#   brew install orchestra
#
# Or in one command:
#   brew install orchestra-mcp/tap/orchestra
#
class Orchestra < Formula
  desc "AI-agentic project management via Model Context Protocol (MCP)"
  homepage "https://github.com/orchestra-mcp/framework"
  version "VERSION_PLACEHOLDER"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/orchestra-mcp/framework/releases/download/v#{version}/orchestra-darwin-arm64.tar.gz"
      sha256 "SHA256_DARWIN_ARM64"
    end

    on_intel do
      url "https://github.com/orchestra-mcp/framework/releases/download/v#{version}/orchestra-darwin-amd64.tar.gz"
      sha256 "SHA256_DARWIN_AMD64"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/orchestra-mcp/framework/releases/download/v#{version}/orchestra-linux-arm64.tar.gz"
      sha256 "SHA256_LINUX_ARM64"
    end

    on_intel do
      url "https://github.com/orchestra-mcp/framework/releases/download/v#{version}/orchestra-linux-amd64.tar.gz"
      sha256 "SHA256_LINUX_AMD64"
    end
  end

  def install
    bin.install "orchestra"
    bin.install "orchestrator"
    bin.install "storage-markdown"
    bin.install "tools-features"
    bin.install "transport-stdio"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/orchestra version")
  end
end
