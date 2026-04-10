#!/usr/bin/env node
/**
 * Thin wrapper that locates the yummycli binary installed by install.js
 * and forwards all arguments, stdin, stdout, and stderr directly to it.
 */

"use strict";

const path = require("path");
const { execFileSync } = require("child_process");

const isWindows = process.platform === "win32";
const binaryName = isWindows ? "yummycli.exe" : "yummycli";
const binaryPath = path.join(__dirname, "..", "bin", binaryName);

try {
  execFileSync(binaryPath, process.argv.slice(2), { stdio: "inherit" });
} catch (err) {
  process.exit(err.status ?? 1);
}
