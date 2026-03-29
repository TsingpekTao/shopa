#!/usr/bin/env node
const { spawn } = require("child_process");
const fs = require("fs");
const path = require("path");

const appName = process.argv[2];
const action = process.argv[3] || "dev";

if (!appName) {
  console.error("[shopa-web] missing app name, e.g. node ./scripts/run-next-app.cjs mall-web dev");
  process.exit(1);
}

const appDir = path.resolve(__dirname, `../apps/${appName}`);
if (!fs.existsSync(path.join(appDir, "package.json"))) {
  console.error(`[shopa-web] app package not found: ${appDir}`);
  process.exit(1);
}

let nextCli = "";
try {
  nextCli = require.resolve("next/dist/bin/next", {
    paths: [appDir, path.resolve(__dirname, "..")]
  });
} catch (_err) {
  nextCli = "";
}
if (!nextCli || !fs.existsSync(nextCli)) {
  console.error(`[shopa-web] Next.js not found for ${appName}. Please run npm install first.`);
  process.exit(1);
}

const supported = new Set(["dev", "build", "lint"]);
if (!supported.has(action)) {
  console.error(`[shopa-web] unsupported action: ${action}`);
  process.exit(1);
}

process.env.NEXT_IGNORE_INCORRECT_LOCKFILE = process.env.NEXT_IGNORE_INCORRECT_LOCKFILE || "1";
if (!process.env.NODE_PATH) {
  const appNodeModules = path.join(appDir, "node_modules");
  const rootNodeModules = path.resolve(__dirname, "../node_modules");
  process.env.NODE_PATH = [appNodeModules, rootNodeModules].join(path.delimiter);
}

const args = [nextCli, action];
if (action === "dev" && process.env.PORT) {
  args.push("-p", process.env.PORT);
}

const child = spawn(process.execPath, args, {
  cwd: appDir,
  stdio: "inherit",
  env: process.env
});

child.on("exit", (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal);
    return;
  }
  process.exit(code ?? 0);
});
