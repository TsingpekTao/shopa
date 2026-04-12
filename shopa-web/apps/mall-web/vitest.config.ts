import { webcrypto } from "crypto";
import path from "path";
import { defineConfig } from "vitest/config";

if (!globalThis.crypto?.getRandomValues && webcrypto?.getRandomValues) {
  globalThis.crypto = webcrypto;
}

export default defineConfig({
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  test: {
    include: [
      "src/features/refund/helpers.test.ts",
      "src/features/agent/api.test.ts",
      "src/features/agent/assistant-display.test.ts",
      "src/features/agent/assistant-copy.test.ts",
      "src/features/agent/assistant-order-selection-card.test.tsx",
    ],
    environment: "node",
    globals: true,
    setupFiles: "./vitest.setup.ts"
  }
});
