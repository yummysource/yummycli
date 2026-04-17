#!/usr/bin/env node
/**
 * postinstall script — downloads the yummycli binary for the current platform
 * from GitHub Releases and places it in the package's bin/ directory.
 *
 * Naming convention (must match .goreleaser.yml archives.name_template):
 *   yummycli-<version>-<os>-<arch>.tar.gz   (non-Windows)
 *   yummycli-<version>-<os>-<arch>.zip      (Windows)
 */

"use strict";

const https = require("https");
const fs = require("fs");
const path = require("path");
const os = require("os");
const { execSync } = require("child_process");

// ── Platform mapping ────────────────────────────────────────────────────────

const PLATFORM_MAP = { darwin: "darwin", linux: "linux", win32: "windows" };
const ARCH_MAP = { x64: "amd64", arm64: "arm64" };

const platform = PLATFORM_MAP[process.platform];
const arch = ARCH_MAP[process.arch];

if (!platform || !arch) {
  console.error(
    `yummycli: unsupported platform ${process.platform}/${process.arch}`
  );
  process.exit(1);
}

// ── Paths ───────────────────────────────────────────────────────────────────

const pkg = require("../package.json");
const version = pkg.version;
const isWindows = platform === "windows";
const archiveName = isWindows
  ? `yummycli-${version}-${platform}-${arch}.zip`
  : `yummycli-${version}-${platform}-${arch}.tar.gz`;

const downloadURL = `https://github.com/yummysource/yummycli/releases/download/v${version}/${archiveName}`;
const binDir = path.join(__dirname, "..", "bin");
const binaryName = isWindows ? "yummycli.exe" : "yummycli";
const binaryDest = path.join(binDir, binaryName);

// ── Helpers ──────────────────────────────────────────────────────────────────

/**
 * Follow HTTP redirects and resolve to the final response.
 * @param {string} url
 * @returns {Promise<import("http").IncomingMessage>}
 */
function get(url) {
  return new Promise((resolve, reject) => {
    https
      .get(url, (res) => {
        if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
          resolve(get(res.headers.location));
        } else {
          resolve(res);
        }
      })
      .on("error", reject);
  });
}

/**
 * Download url to destPath.
 * @param {string} url
 * @param {string} destPath
 * @returns {Promise<void>}
 */
function download(url, destPath) {
  return new Promise(async (resolve, reject) => {
    const res = await get(url);
    if (res.statusCode !== 200) {
      reject(new Error(`HTTP ${res.statusCode} downloading ${url}`));
      return;
    }
    const out = fs.createWriteStream(destPath);
    res.pipe(out);
    out.on("finish", resolve);
    out.on("error", reject);
  });
}

// ── Helpers ── skills ────────────────────────────────────────────────────────

/**
 * Attempt to update agent skills via the `skills` CLI.
 * This is best-effort: a missing `skills` binary is not an error.
 */
function installSkills() {
  // Determine whether `skills` is available on PATH before attempting npx,
  // so we don't silently download an unrelated package.
  try {
    execSync("skills --version", { stdio: "ignore" });
  } catch {
    console.log(
      "yummycli: agent skills not updated — run manually:\n" +
        "  npx skills add yummysource/yummycli -y -g"
    );
    return;
  }

  try {
    console.log("yummycli: updating agent skills...");
    execSync("skills add yummysource/yummycli -y -g", { stdio: "inherit" });
    console.log("yummycli: agent skills updated");
  } catch {
    console.log(
      "yummycli: agent skills update failed — run manually:\n" +
        "  npx skills add yummysource/yummycli -y -g"
    );
  }
}

// ── Main ─────────────────────────────────────────────────────────────────────

async function main() {
  fs.mkdirSync(binDir, { recursive: true });

  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), "yummycli-"));
  const archivePath = path.join(tmpDir, archiveName);

  console.log(`yummycli: downloading ${downloadURL}`);

  try {
    await download(downloadURL, archivePath);

    if (isWindows) {
      execSync(
        `powershell -Command "Expand-Archive -Path '${archivePath}' -DestinationPath '${tmpDir}' -Force"`
      );
    } else {
      execSync(`tar -xzf "${archivePath}" -C "${tmpDir}"`);
    }

    const extracted = path.join(tmpDir, binaryName);
    fs.copyFileSync(extracted, binaryDest);

    if (!isWindows) {
      fs.chmodSync(binaryDest, 0o755);
    }

    console.log(`yummycli: installed to ${binaryDest}`);
  } finally {
    fs.rmSync(tmpDir, { recursive: true, force: true });
  }

  installSkills();
}

main().catch((err) => {
  console.error("yummycli: install failed:", err.message);
  process.exit(1);
});
