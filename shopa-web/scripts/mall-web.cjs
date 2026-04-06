#!/usr/bin/env node
const { spawn } = require("child_process");
const path = require("path");

const action = process.argv[2] || "dev";
if (action === "dev" && !process.env.PORT) {
  process.env.PORT = "3100";
}
const child = spawn(process.execPath, [path.resolve(__dirname, "./run-next-app.cjs"), "mall-web", action], {
  stdio: "inherit",
  env: process.env
});

child.on("exit", (code) => process.exit(code ?? 0));
