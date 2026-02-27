#!/usr/bin/env node
//
// Postinstall script for @orchestra-mcp/cli npm package.
// Downloads the correct Go binary tarball for the current platform.
//
const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");
const https = require("https");

const REPO = "orchestra-mcp/framework";
const VERSION = require("../package.json").version;
const BIN_DIR = path.join(__dirname, "..", "bin");

const PLATFORM_MAP = {
  darwin: "darwin",
  linux: "linux",
};

const ARCH_MAP = {
  x64: "amd64",
  arm64: "arm64",
};

function getPlatform() {
  const os = PLATFORM_MAP[process.platform];
  const arch = ARCH_MAP[process.arch];

  if (!os || !arch) {
    console.error(
      `Unsupported platform: ${process.platform}/${process.arch}`
    );
    console.error("Orchestra supports: darwin/linux on amd64/arm64");
    process.exit(1);
  }

  return `${os}-${arch}`;
}

function download(url, dest) {
  return new Promise((resolve, reject) => {
    const follow = (url) => {
      https
        .get(url, (res) => {
          if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
            follow(res.headers.location);
            return;
          }
          if (res.statusCode !== 200) {
            reject(new Error(`Download failed: HTTP ${res.statusCode}`));
            return;
          }
          const file = fs.createWriteStream(dest);
          res.pipe(file);
          file.on("finish", () => {
            file.close(resolve);
          });
        })
        .on("error", reject);
    };
    follow(url);
  });
}

async function main() {
  const platform = getPlatform();
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/orchestra-${platform}.tar.gz`;
  const tarball = path.join(BIN_DIR, "orchestra.tar.gz");

  console.log(`Orchestra MCP: downloading ${platform} binary...`);

  fs.mkdirSync(BIN_DIR, { recursive: true });

  await download(url, tarball);

  console.log("Extracting...");
  execSync(`tar -xzf "${tarball}" -C "${BIN_DIR}"`, { stdio: "inherit" });
  fs.unlinkSync(tarball);

  // Make binaries executable.
  const binaries = [
    "orchestra",
    "orchestrator",
    "storage-markdown",
    "tools-features",
    "transport-stdio",
  ];
  for (const bin of binaries) {
    const binPath = path.join(BIN_DIR, bin);
    if (fs.existsSync(binPath)) {
      fs.chmodSync(binPath, 0o755);
    }
  }

  console.log("Orchestra MCP installed successfully.");
  console.log("Run: npx orchestra init");
}

main().catch((err) => {
  console.error("Installation failed:", err.message);
  console.error("You can install manually: https://github.com/" + REPO);
  process.exit(1);
});
